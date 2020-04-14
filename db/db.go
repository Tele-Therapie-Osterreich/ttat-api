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

// ProfileSelection is an enumerated type used to determine which
// profile (public or pending) is returned from queries.
type ProfileSelection uint

// Enumeration values for ProfileSelection.
const (
	PublicOnly ProfileSelection = iota
	PendingOnly
	PreferPublic
)

// DB describes the database operations.
type DB interface {
	// ----------------------------------------------------------------------
	//
	// LOGIN PROCESSING

	// Login performs login actions for a given email address:
	//
	//  - If an account with the given email address does not already
	//    exist in the database, then create a new user account with the
	//    given email address. Also create a new non-public therapist
	//    profile, defaulting all user information fields to empty.
	//
	//  - If an account with the give email address does exist, return
	//    the therapist's public profile if they have one and their
	//    pending one if not.
	Login(email string) (*model.TherapistInfo, bool, error)

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

	// TherapistByID returns the therapist model for a given therapist
	// ID.
	TherapistByID(thID int) (*model.Therapist, error)

	// TherapistInfoByID returns full therapist information for a given
	// therapist ID.
	TherapistInfoByID(thID int, profile ProfileSelection) (*model.TherapistInfo, error)

	// TherapistInfoByEmail returns full therapist information for a
	// given therapist email address, preferentially returning any
	// public profile.
	TherapistInfoByEmail(email string) (*model.TherapistInfo, error)

	// TODO: ALLOW UPDATES TO THERAPIST EMAIL ADDRESS? EVERYTHING ELSE
	// GOES THROUGH PROFILE EDITS...

	// DeleteTherapist deletes the given therapist account.
	// TODO: CHECK DELETION OF ASSOCIATED DATA -- EVERYTHING LINKED TO
	// IDs IN therapists TABLE SHOULD HAVE "ON DELETE CASCADE".
	DeleteTherapist(thID int) error

	// ----------------------------------------------------------------------
	//
	// THERAPIST PROFILES

	// TherapistProfileByTherapistID retrieves the public or pending
	// profile of a given therapist.
	TherapistProfileByTherapistID(thID int, public bool) (*model.TherapistProfile, error)

	// UpdateTherapistProfile performs profile updating for a therapist.
	// If the therapist already has a pending profile, the update
	// profile replaces the pending profile. If the therapist does not
	// have a pending profile, a new one is generated from the profile
	// passed in.
	UpdateTherapistProfile(thID int, patch []byte) (*model.ImagePatch, error)

	// AbandonTherapistEdits deletes any pending edits profile
	// associated with an active therapist.
	AbandonTherapistEdits(thID int) error

	// ----------------------------------------------------------------------
	//
	// SEARCH

	// TherapistSearch performs the "matchmaker" searching, filtering
	// public therapist profiles by therapist type and a list of
	// specialities. Results are returned as views suitable for display
	// in a summary list.
	//TherapistSearch(t types.TherapistType, specialities []string) ([]*model.TherapistSummaryView, error)

	// ----------------------------------------------------------------------
	//
	// IMAGES

	// Retrieve image.
	ImageByID(imgID int) (*model.Image, error)

	// Retrieve therapist profile image.
	ImageByProfileID(prID int) (*model.Image, error)

	// Insert or update image.
	UpsertImage(image *model.Image) (*model.Image, error)
}

//go:generate go-bindata -pkg db -o migrations.go migrations/...
