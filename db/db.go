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

// ErrTherapistNotFound is the error returned when an attempt is made
// to access or manipulate a therapist with an unknown ID.
var ErrTherapistNotFound = errors.New("therapist ID not found")

// ErrImageNotFound is the error returned when an attempt is made to
// access or manipulate a therapist image with an unknown ID.
var ErrImageNotFound = errors.New("therapist image ID not found")

// ErrReadOnlyField is the error returned when an attempt is made to
// update a read-only field for a therapist (e.g. email, last login
// time).
var ErrReadOnlyField = errors.New("attempt to modify read-only field")

// DB describes the database operations.
type DB interface {
	// ----------------------------------------------------------------------
	//
	// LOGIN PROCESSING

	// Login performs login actions for a given email address:
	//
	// If an account with the given email address does not already exist
	// in the database, then create a new user account with the given
	// email address, defaulting all user information fields to empty.
	//
	// Returns the full therapist record of the logged in therapist.
	Login(email string) (*model.Therapist, *model.Image, bool, error)

	// ----------------------------------------------------------------------
	//
	// MAGIC EMAIL LOGIN PROCESSING

	// CreateLoginToken creates a new unique six-digit numerical login
	// token for the given email. It returns the token directly.
	CreateLoginToken(email string, language string) (string, error)

	// CheckLoginToken checks whether a given login token is valid and
	// has not expired. If the token is good, the email address and
	// language associated with it are returned.
	CheckLoginToken(token string) (string, string, error)

	// ----------------------------------------------------------------------
	//
	// SESSION HANDLING

	// CreateSession generates a new session token for a therapist (or
	// reconnects to an existing session for the given therapist).
	CreateSession(thID int) (string, error)

	// LookupSession checks a session token and returns the associated
	// therapist ID if the session is known.
	LookupSession(token string) (*int, error)

	// DeleteSessions deletes all sessions for a therapist, i.e. logs
	// the therapist out of all devices where they're logged in.
	DeleteSessions(thID int) error

	// ----------------------------------------------------------------------
	//
	// THERAPISTS

	// TherapistByID returns the full therapist model for a given
	// therapist ID.
	TherapistByID(thID int) (*model.Therapist, error)

	// UpdateTherapist updates the therapist's details in the database.
	// The id, email, status and created_at fields are read-only using
	// this method.
	UpdateTherapist(therapist *model.Therapist) error

	// DeleteTherapist deletes the given therapist account.
	DeleteTherapist(thID int) error

	// ----------------------------------------------------------------------
	//
	// IMAGES

	// Retrieve image.
	ImageByID(imgID int) (*model.Image, error)

	// Retrieve therapist image.
	ImageByTherapistID(thID int) (*model.Image, error)

	// Insert or update image.
	UpsertImage(image *model.Image) (*model.Image, error)
}

//go:generate go-bindata -pkg db -o migrations.go migrations/...
