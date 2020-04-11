package server

import (
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/Tele-Therapie-Osterreich/ttat-api/messages"
	"github.com/Tele-Therapie-Osterreich/ttat-api/model"
)

func (s *Server) RequestLoginEmail(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// Decode request body.
	body := messages.ReqLoginEmailRequest{}
	err := Unmarshal(r.Body, &body)
	if err != nil {
		return BadRequest(w, err.Error())
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

func (s *Server) Login(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// Decode request body to get login token.
	req := messages.LoginRequest{}
	err := Unmarshal(r.Body, &req)
	if err != nil {
		return BadRequest(w, err.Error())
	}

	// Look up login token and if not found or expired, return error.
	email, _, err := s.db.CheckLoginToken(req.LoginToken)
	if err != nil {
		return BadRequest(w, "Unknown login token")
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

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// Get session cookie. If there is no session, this is a no-op.
	authInfo := AuthInfoFromContext(r.Context())
	if !authInfo.Authenticated {
		return NotFound(w)
	}

	if authInfo.UserID != 0 {
		// Delete all sessions from database for matching user ID.
		s.db.DeleteSessions(authInfo.UserID)
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
