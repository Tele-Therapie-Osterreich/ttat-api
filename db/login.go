package db

import (
	"database/sql"

	"github.com/Tele-Therapie-Osterreich/ttat-api/model"
)

// Login performs login actions for a given email address:
//
// If an account with the given email address does not already exist
// in the database, then create a new user account with the given
// email address, defaulting all user information fields to empty.
//
// Returns the full therapist record of the logged in therapist.
func (pg *PGClient) Login(email string) (*model.Therapist, *model.Image, bool, error) {
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

	th := &model.Therapist{}
	err = tx.Get(th, therapistByEmail, email)
	if err == nil {
		image, err := pg.ImageByTherapistID(th.ID)
		if err != ErrImageNotFound && err != nil {
			return nil, nil, false, err
		}
		return th, image, false, nil
	}

	th = &model.Therapist{
		Email: email,
	}
	rows, err := tx.NamedQuery(createTherapist, th)
	if err != nil {
		return nil, nil, false, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, nil, false, sql.ErrNoRows
	}

	err = rows.Scan(&th.ID)
	if err != nil {
		return nil, nil, false, err
	}

	return th, nil, true, nil
}

const therapistByEmail = `
SELECT id, email, type, name,
       street_address, city, postcode, country,
       phone, languages, short_profile, full_profile,
       status, created_at
  FROM therapists
 WHERE email = $1`

const createTherapist = `
INSERT INTO therapists (email) VALUES (:email) RETURNING id`
