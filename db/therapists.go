package db

import (
	"database/sql"

	"github.com/Tele-Therapie-Osterreich/ttat-api/model"
	"github.com/jmoiron/sqlx"

	// Import Postgres DB driver.
	_ "github.com/lib/pq"
)

// TherapistByID looks up a therapist by their therapist ID.
func (pg *PGClient) TherapistByID(thID int) (*model.Therapist, error) {
	return therapistByID(pg.DB, thID)
}

func therapistByID(q sqlx.Queryer, id int) (*model.Therapist, error) {
	th := &model.Therapist{}
	err := sqlx.Get(q, th, qTherapistByID, id)
	if err == sql.ErrNoRows {
		return nil, ErrTherapistNotFound
	}
	if err != nil {
		return nil, err
	}
	return th, nil
}

const qTherapistByID = `
SELECT id, email, status, last_login_at, created_at
  FROM therapists
 WHERE id = $1`

func therapistByEmail(q sqlx.Queryer, email string) (*model.Therapist, error) {
	th := &model.Therapist{}
	err := sqlx.Get(q, th, qTherapistByEmail, email)
	if err == sql.ErrNoRows {
		return nil, ErrTherapistNotFound
	}
	if err != nil {
		return nil, err
	}
	return th, nil
}

const qTherapistByEmail = `
SELECT id, email, status, last_login_at, created_at
  FROM therapists
 WHERE email = $1`

// DeleteTherapist deletes the given therapist account.
func (pg *PGClient) DeleteTherapist(thID int) error {
	result, err := pg.DB.Exec(qDeleteTherapist, thID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return ErrTherapistNotFound
	}
	return nil
}

const qDeleteTherapist = `DELETE FROM therapists WHERE id = $1`
