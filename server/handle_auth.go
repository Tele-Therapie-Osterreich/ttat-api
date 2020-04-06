package server

import (
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/Tele-Therapie-Osterreich/ttat-api/chassis"
	"github.com/Tele-Therapie-Osterreich/ttat-api/messages"
	"github.com/Tele-Therapie-Osterreich/ttat-api/model"
)

func (s *Server) requestLoginEmail(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// Decode request body.
	body := messages.ReqLoginEmailRequest{}
	err := chassis.Unmarshal(r.Body, &body)
	if err != nil {
		return chassis.BadRequest(w, err.Error())
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
	s.SendEmail("login-email-request", body.Email, body.Language,
		map[string]string{"login_token": token})

	return chassis.NoContent(w)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// Decode request body to get login token.
	req := messages.LoginRequest{}
	err := chassis.Unmarshal(r.Body, &req)
	if err != nil {
		return chassis.BadRequest(w, err.Error())
	}

	// Look up login token and if not found or expired, return error.
	email, _, err := s.db.CheckLoginToken(req.LoginToken)
	if err != nil {
		return chassis.BadRequest(w, "Unknown login token")
	}

	// Perform login processing.
	user, avatar, new, err := s.db.LoginUser(email)
	if err != nil {
		return nil, errors.Wrap(err, "performing login processing")
	}

	// Create a session for the user (or reconnect to an existing
	// session).
	token, err := s.db.CreateSession(user.ID)
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
		UserProfile: model.UserFullProfileFromUser(user, avatar),
		NewUser:     new,
	}
	return &resp, nil
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// Get session cookie. If there is no session, this is a no-op.
	authInfo := AuthInfoFromContext(r.Context())
	if !authInfo.Authenticated {
		return chassis.NotFound(w)
	}

	if authInfo.UserID != 0 {
		// Delete all sessions from database for matching user ID.
		s.db.DeleteSessions(authInfo.UserID)
	}

	s.clearSessionCookie(w)
	return chassis.NoContent(w)
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
