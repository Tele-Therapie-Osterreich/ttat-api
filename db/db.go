package db

import (
	"errors"
	"time"

	"github.com/Tele-Therapie-Osterreich/ttat-api/model"
)

// LoginTokenDuration is the time for which a login token is valid:
// tokens are considered to have expired if they have not been
// presented within this time of the time they are created.
const LoginTokenDuration = 1 * time.Hour

// ErrLoginTokenNotFound is the error returned when an unknown login
// token is submitted for checking.
var ErrLoginTokenNotFound = errors.New("login token not found")

// ErrSessionNotFound is the error returned when an unknown session ID
// is used.
var ErrSessionNotFound = errors.New("session not found")

// ErrUserNotFound is the error returned when an attempt is made
// to access or manipulate a user with an unknown ID.
var ErrUserNotFound = errors.New("user ID not found")

// ErrImageNotFound is the error returned when an attempt is made to
// access or manipulate a user image with an unknown ID.
var ErrImageNotFound = errors.New("user image ID not found")

// ErrReadOnlyField is the error returned when an attempt is made to
// update a read-only field for a user (e.g. email, last login
// time).
var ErrReadOnlyField = errors.New("attempt to modify read-only field")

// DB describes the database operations.
type DB interface {
	// ----------------------------------------------------------------------
	//
	// LOGIN PROCESSING

	// LoginUser performs login actions for a given email address:
	//
	//  - If an account with the given email address already exists in the
	//    database, then set the account's last_login field to the current time.
	//
	//  - If an account with the given email address does not already
	//    exist in the database, then create a new user account with the
	//    given email address, defaulting all user information fields to
	//    empty and setting the new account's last_login field to the
	//    current time.
	//
	// In both cases, return the full user record of the logged in user.
	LoginUser(email string) (*model.User, *model.Image, bool, error)

	// ----------------------------------------------------------------------
	//
	// MAGIC EMAIL LOGIN PROCESSING

	// CreateLoginToken creates a new unique six-digit numerical login
	// token for the given email. It returns the token directly.
	CreateLoginToken(email string, language string) (string, error)

	// CheckLoginToken checks whether a given login token is valid and
	// has not expired. If the token is good, the email address
	// and language associated with it are returned.
	CheckLoginToken(token string) (string, string, error)

	// ----------------------------------------------------------------------
	//
	// SESSION HANDLING

	// CreateSession generates a new session token for a user (or
	// reconnects to an existing session for the given user).
	CreateSession(userID int) (string, error)

	// LookupSession checks a session token and returns the associated
	// user ID if the session is known.
	LookupSession(token string) (*int, error)

	// DeleteSessions deletes all sessions for a user, i.e. logs the
	// user out of all devices where they're logged in.
	DeleteSessions(userID int) error

	// ----------------------------------------------------------------------
	//
	// USERS

	// UserByID returns the full user model for a given user ID.
	UserByID(id int) (*model.User, error)

	// UpdateUser updates the user's details in the database. The id,
	// email, last_login and api_key fields are read-only using this
	// method.
	UpdateUser(therapist *model.User) error

	// DeleteUser deletes the given user account.
	DeleteUser(id int) error

	// Retrieve image.
	ImageByID(id int) (*model.Image, error)

	// Retrieve user image.
	ImageByUserID(id int) (*model.Image, error)

	// Insert or update image.
	UpsertImage(image *model.Image) (*model.Image, error)
}

//go:generate go-bindata -pkg db -o migrations.go migrations/...
