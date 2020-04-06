package db

import (
	"database/sql"
	"time"

	"github.com/Tele-Therapie-Osterreich/ttat-api/chassis"
	"github.com/Tele-Therapie-Osterreich/ttat-api/model"

	// Import Postgres DB driver.
	_ "github.com/lib/pq"
)

// CreateLoginToken creates a new unique six-digit numerical login
// token for the given email. It returns the token directly.
func (pg *PGClient) CreateLoginToken(email string, language string) (string, error) {
	// We deliberately don't validate the email address here. Email
	// addresses are a mess, and the simplest thing to do is just to
	// send an email to the address. If it doesn't work, then the entry
	// in the login_tokens table we create here will get cleaned up
	// after the token expires, and no harm done.

	// Do this in a transaction, since we want to clear out existing
	// tokens for the email address and ensure that we respect the token
	// uniqueness constraint.
	tx, err := pg.DB.Beginx()
	if err != nil {
		return "", err
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	// Delete any existing tokens for this email, and at the same time
	// delete any expired tokens.
	_, err = tx.Exec(deleteTokenForEmail, email)
	if err != nil {
		return "", err
	}

	// Create and insert a unique token.
	token := RandToken()
	for {
		result, err := tx.Exec(insertToken,
			token, email, language, time.Now().Add(LoginTokenDuration))
		if err != nil {
			return "", err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return "", err
		}
		if rows == 1 {
			break
		}

		// Token collision: try again...
		token = RandToken()
	}

	return token, nil
}

const deleteTokenForEmail = `
DELETE FROM login_tokens
 WHERE email = $1 OR expires_at < NOW()`

const insertToken = `
INSERT INTO login_tokens (token, email, language, expires_at)
    VALUES ($1, $2, $3, $4)
    ON CONFLICT DO NOTHING`

// CheckLoginToken checks whether a given login token is valid and has
// not expired. If the token is good, the email address associated
// with it is returned.
func (pg *PGClient) CheckLoginToken(token string) (string, string, error) {
	// In a transaction...
	tx, err := pg.DB.Beginx()
	if err != nil {
		return "", "", err
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	// Look up token.
	tokInfo := TokenInfo{}
	err = tx.Get(&tokInfo, lookupToken, token)
	if err == sql.ErrNoRows {
		return "", "", ErrLoginTokenNotFound
	}
	if err != nil {
		return "", "", err
	}

	// Clear token entries.
	tx.Exec(cleanupTokens, token)

	return tokInfo.Email, tokInfo.Language, nil
}

// TokenInfo is a temporary structure for holding information about
// the login token being retrieved.
type TokenInfo struct {
	Email    string `db:"email"`
	Site     string `db:"site"`
	Language string `db:"language"`
}

const lookupToken = `
SELECT email, language FROM login_tokens
 WHERE token = $1 AND expires_at >= NOW()`

const cleanupTokens = `
DELETE FROM login_tokens WHERE token = $1 OR expires_at < NOW()`

// CreateSession generates a new session token for a therapist user,
// or reconnects to an existing session for the requesting user ID.
func (pg *PGClient) CreateSession(userID int) (string, error) {
	// In a transaction...
	tx, err := pg.DB.Beginx()
	if err != nil {
		return "", err
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	// Look up existing session.
	token := ""
	err = tx.Get(&token, sessionByUser, userID)
	if err == nil {
		return token, nil
	}
	if err != sql.ErrNoRows {
		return "", err
	}

	// Create new session.
	id := chassis.NewID(16)

	_, err = tx.Exec(`INSERT INTO sessions VALUES ($1, $2)`, id, userID)
	if err != nil {
		return "", err
	}

	return id, nil
}

const sessionByUser = `
SELECT token FROM sessions WHERE user_id = $1`

// LookupSession checks a session token and returns the associated
// user ID, email and admin flag if the session is known.
func (pg *PGClient) LookupSession(token string) (*int, error) {
	sess := model.Session{}
	err := pg.DB.Get(&sess, lookupSession, token)
	if err == sql.ErrNoRows {
		return nil, ErrSessionNotFound
	}
	if err != nil {
		return nil, err
	}
	return sess.UserID, nil
}

const lookupSession = `
SELECT token, user_id
  FROM sessions
 WHERE token = $1`

// DeleteSessions deletes all sessions for a user, i.e. logs the user
// out of all devices where they're logged in.
func (pg *PGClient) DeleteSessions(userID int) error {
	_, err := pg.DB.Exec(`DELETE FROM sessions WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}
	return nil
}

// Parameters to use for random token generation.
const (
	tokenLen     = 6
	dChars       = "0123456789"
	dCharIdxBits = 4                   // 4 bits to represent a character index
	dCharIdxMask = 1<<dCharIdxBits - 1 // All 1-bits, as many as dCharIdxBits
	dCharIdxMax  = 63 / dCharIdxBits   // # of char indices fitting in 63 bits
)

// RandToken generates a random string of digits to use as a login token.
func RandToken() string {
	return chassis.RandString(dChars, dCharIdxBits, dCharIdxMask, dCharIdxMax,
		tokenLen)
}
