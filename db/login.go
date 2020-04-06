package db

import (
	"database/sql"

	"github.com/rs/zerolog/log"

	"github.com/Tele-Therapie-Osterreich/ttat-api/model"
)

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
func (pg *PGClient) LoginUser(email string) (*model.User, *model.Image, bool, error) {
	tx, err := pg.DB.Beginx()
	if err != nil {
		return nil, nil, false, err
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	user := &model.User{}
	err = tx.Get(user, userByEmail, email)
	if err == nil {
		image, err := pg.ImageByUserID(user.ID)
		if err != ErrImageNotFound && err != nil {
			return nil, nil, false, err
		}
		return user, image, false, nil
	}

	log.Info().Msg("NEW USER")
	user = &model.User{
		Email: email,
	}
	rows, err := tx.NamedQuery(createUser, user)
	if err != nil {
		return nil, nil, false, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, nil, false, sql.ErrNoRows
	}

	err = rows.Scan(&user.ID)
	if err != nil {
		return nil, nil, false, err
	}

	return user, nil, true, nil
}

const userByEmail = `
SELECT id, email, type, name,
       street_address, city, postcode, country,
       phone, short_profile, full_profile,
       status, created_at
  FROM users
 WHERE email = $1`

const createUser = `
INSERT INTO users (email) VALUES (:email) RETURNING id`
