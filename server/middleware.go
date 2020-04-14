package server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

// Add middleware specific to API gateway.
func (s *Server) addMiddleware(r chi.Router, devMode bool,
	csrfSecret string, corsOrigins []string) {
	// IP header middleware.
	r.Use(middleware.RealIP)

	// Add common middleware.
	addBasicMiddleware(r)

	// Very basic concurrent request service throttling.
	// TODO: PER-ROUTE, PER-IP RATE LIMITING.
	r.Use(middleware.Throttle(1000))

	// CORS configuration.
	opts := cors.Options{
		AllowedOrigins:   corsOrigins,
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "PUT", "PATCH", "GET", "DELETE"},
		// TODO: SORT THIS NEXT SETTING OUT...
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"X-CSRF-Token"},
	}
	co := cors.New(opts)
	r.Use(co.Handler)

	// Add middleware to implement double token submission CSRF
	// protection. One token is stored in a secure HTTP-only cookie
	// (security switched off for local development...), and the other is
	// sent in each response in an X-CSRF-Token header. This second token
	// must be submitted in each request that makes state changes by
	// adding an X-CSRF-Token header to the request.
	options := []csrf.Option{}
	options = append(options, csrf.CookieName("csrftoken"))
	if devMode {
		options = append(options, csrf.Secure(false))
	}
	csrfMiddleware := csrf.Protect([]byte(csrfSecret), options...)
	r.Use(csrfMiddleware)

	// Extra middleware to add CSRF token to all responses.
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-CSRF-Token", csrf.Token(r))
			next.ServeHTTP(w, r)
		})
	})
}

// AddBasicMiddleware adds basic middleware for all routes.
func addBasicMiddleware(r chi.Router) {
	// Set up zerolog request logging.
	r.Use(hlog.NewHandler(log.Logger))
	logs := func(r *http.Request, status, size int, duration time.Duration) {
		basicRequestLog(r, status, size, duration).Msg("")
	}
	r.Use(hlog.AccessHandler(logs))
	r.Use(hlog.RemoteAddrHandler("ip"))
	r.Use(hlog.UserAgentHandler("user_agent"))
	r.Use(hlog.RefererHandler("referer"))

	// TODO: MAKE THIS LESS VERBOSE FOR PRODUCTION.
	r.Use(logHandler)

	// Panic recovery.
	r.Use(middleware.Recoverer)

	// QUESTION: I DON'T THINK THIS WILL WORK WITHOUT MAKING ALL THE
	// HANDLERS CONTEXT-AWARE. WHAT'S THE RIGHT WAY TO DEAL WITH THIS?
	// r.Use(middleware.Timeout(60 * time.Second))
}

func logHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" {
			next.ServeHTTP(w, r)
			return
		}

		x, err := httputil.DumpRequest(r, false)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		save := r.Body
		if r.Body == nil {
			r.Body = nil
		} else {
			save, r.Body, _ = drainBody(r.Body)
		}
		reqBody, err := ioutil.ReadAll(save)
		if !utf8.Valid(reqBody) {
			reqBody = []byte("** BINARY DATA IN BODY **")
		}
		x = append(x, reqBody...)
		log.Info().Str("dir", "request").Msg(string(x))
		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)
		resp := fmt.Sprintf("%d\n", rec.Code)
		for k, v := range rec.HeaderMap {
			resp += k + ": " + strings.Join(v, ",") + "\n"
		}
		body := rec.Body.String()
		if !utf8.Valid([]byte(body)) {
			body = "** BINARY DATA IN BODY **"
		}
		log.Info().Str("dir", "response").Msg(resp + body)

		// this copies the recorded response to the response writer
		for k, v := range rec.HeaderMap {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Code)
		rec.Body.WriteTo(w)
	})
}

func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

// Basic HTTP request logging.
func basicRequestLog(r *http.Request, status, size int, duration time.Duration) *zerolog.Event {
	if r.URL.Path == "/healthz" {
		return nil
	}
	return hlog.FromRequest(r).Info().
		Str("method", r.Method).
		Str("url", r.URL.String()).
		Int("status", status).
		Int("size", size).
		Dur("duration", duration)
}

// AuthInfo carries authentication information derived from the
// session. This provides enough information to do authentication and
// authorisation checking.
type AuthInfo struct {
	// Is the request authenticated?
	Authenticated bool

	// User ID for the authenticated therapist user making the request
	// (will be zero for an unauthenticated request or a request from a
	// patient).
	TherapistID int
}

// ctxKey is a key type for request context information.
type ctxKey int

// Request context keys that we use.
const (
	authInfoCtxKey ctxKey = iota
)

// CredentialCtx extracts credential information from the request
// (either via a session cookie or an API key), looks up the
// corresponding user information and injects the resulting
// authorisation information into the request context as a
// chassis.AuthInfo value.
func CredentialCtx(s *Server) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get session cookie.
			session := ""
			if cookie, err := r.Cookie("session"); err == nil {
				session = cookie.Value
			}
			if session == "" {
				next.ServeHTTP(w, r)
				return
			}

			authInfo := AuthInfo{}
			thID, err := s.db.LookupSession(session)
			if err == nil {
				authInfo.Authenticated = true
				authInfo.TherapistID = *thID
			}
			next.ServeHTTP(w, r.WithContext(NewAuthContext(r.Context(), &authInfo)))
		})
	}
}

// NewAuthContext returns a new Context that carries authentication
// information.
func NewAuthContext(ctx context.Context, info *AuthInfo) context.Context {
	return context.WithValue(ctx, authInfoCtxKey, info)
}

// AuthInfoFromContext returns the authentication information in a
// context, if any.
func AuthInfoFromContext(ctx context.Context) *AuthInfo {
	info, ok := ctx.Value(authInfoCtxKey).(*AuthInfo)
	if !ok {
		info = &AuthInfo{}
	}
	return info
}
