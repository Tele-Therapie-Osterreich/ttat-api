package server

import (
	"net/http"
	"regexp"
	"time"

	"github.com/pkg/errors"

	"github.com/Tele-Therapie-Osterreich/ttat-api/messages"
	"github.com/Tele-Therapie-Osterreich/ttat-api/model"
)

var emailRE = regexp.MustCompile(`(?i)^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$`)

// RequestLoginEmail handles requests to /auth/request-login-email.
func (s *Server) requestLoginEmail(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// Decode request body.
	body := messages.ReqLoginEmailRequest{}
	err := Unmarshal(r.Body, &body)
	if err != nil {
		return BadRequest(w, messages.APIError{
			Code:    messages.InvalidJSONBody,
			Message: "invalid JSON body for login email request",
		})
	}

	// Check email address is reasonable.
	if !emailRE.MatchString(body.Email) {
		return BadRequest(w, messages.APIError{
			Code:    messages.InvalidEmailAddress,
			Message: "invalid email address for login email request",
			Field:   body.Email,
		})
	}

	// Default email language to English.
	if body.Language == "" {
		body.Language = "en"
	}

	// Create login token and send email.
	token, err := s.db.CreateLoginToken(body.Email, body.Language)
	if err != nil {
		return nil, err
	}
	s.mailer.Send("login-email-request", body.Email, body.Language,
		map[string]string{"login_token": token})

	return NoContent(w)
}

// Login handles requests to /auth/login.
func (s *Server) login(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// Decode request body to get login token.
	req := messages.LoginRequest{}
	err := Unmarshal(r.Body, &req)
	if err != nil {
		return BadRequest(w, messages.APIError{
			Code:    messages.InvalidJSONBody,
			Message: "invalid JSON body for login request",
		})
	}

	// Look up login token and if not found or expired, return error.
	email, _, err := s.db.CheckLoginToken(req.LoginToken)
	if err != nil {
		return BadRequest(w, messages.APIError{
			Code:    messages.UnknownLoginToken,
			Message: "unknown login token for login request",
		})
	}

	// Perform login processing.
	info, new, err := s.db.Login(email)
	if err != nil {
		return nil, errors.Wrap(err, "performing login processing")
	}

	// Create a session for the user (or reconnect to an existing
	// session).
	token, err := s.db.CreateSession(info.Base.ID)
	if err != nil {
		return nil, err
	}

	// Set the session cookie and return the user information as a JSON
	// response.
	auth := http.Cookie{
		Name:     "session",
		Value:    token,
		Secure:   s.secureSession,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &auth)

	// Return user response for marshalling.
	resp := messages.LoginResponse{
		Profile:      model.NewTherapistLoginView(info),
		NewTherapist: new,
	}
	return &resp, nil
}

// logout handles requests to /auth/logout.
func (s *Server) logout(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// Get session cookie. If there is no session, this is a no-op.
	authInfo := AuthInfoFromContext(r.Context())
	if !authInfo.Authenticated {
		return NotFound(w)
	}

	if authInfo.TherapistID != 0 {
		// Delete all sessions from database for matching user ID.
		s.db.DeleteSessions(authInfo.TherapistID)
	}

	s.clearSessionCookie(w)
	return NoContent(w)
}

// Delete session cookie by setting expiry in past.
func (s *Server) clearSessionCookie(w http.ResponseWriter) {
	delAuth := http.Cookie{
		Name:     "session",
		Value:    "",
		Secure:   s.secureSession,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
	}
	http.SetCookie(w, &delAuth)
}
