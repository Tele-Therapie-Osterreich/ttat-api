package model

import "time"

// LoginToken holds information associating a unique six-digit login
// token with a therapist user email. Login tokens have an expiry time
// beyond which they are no longer considered valid.
type LoginToken struct {
	// Token string.
	Token string `db:"token"`

	// Email associated with this token.
	Email string `db:"email"`

	// Language for email sending.
	Language string `db:"language"`

	// Token expiry time.
	ExpiresAt time.Time `db:"expires_at"`
}

// Session holds the information that associates session cookies with
// users. Sessions are long-lived, and are only deleted when a user
// logs out.
type Session struct {
	// Session token.
	Token string `db:"token"`

	// ID of a therapist user the session is associated with.
	UserID *int `db:"user_id"`
}
