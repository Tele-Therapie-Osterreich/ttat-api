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
func (pg *PGClient) Login(email string) (*model.TherapistInfo, bool, error) {
	tx, err := pg.DB.Beginx()
	if err != nil {
		return nil, false, err
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	// TODO: DEAL WITH THERAPIST IMAGES
	info, err := therapistInfoByEmail(tx, email)
	if err == nil {
		// TODO: UPDATE last_login_at HERE
		return info, false, nil
	}
	if err != nil && err != ErrTherapistNotFound {
		return nil, false, err
	}

	th := &model.Therapist{
		Email: email,
	}
	rows, err := tx.NamedQuery(qCreateTherapist, th)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, false, sql.ErrNoRows
	}
	err = rows.Scan(&th.ID)
	if err != nil {
		return nil, false, err
	}
	rows.Close()

	p := &model.TherapistProfile{
		TherapistID: th.ID,
		Public:      false,
	}
	rows2, err := tx.NamedQuery(qCreateTherapistProfile, p)
	if err != nil {
		return nil, false, err
	}
	defer rows2.Close()
	if !rows2.Next() {
		return nil, false, sql.ErrNoRows
	}
	err = rows2.Scan(&p.ID)
	if err != nil {
		return nil, false, err
	}

	return &model.TherapistInfo{
		Base:              th,
		Profile:           p,
		Image:             nil,
		HasPublicProfile:  false,
		HasPendingProfile: true,
	}, true, nil
}

const qCreateTherapist = `
INSERT INTO therapists (email) VALUES (:email) RETURNING id`

const qCreateTherapistProfile = `
INSERT INTO profiles
  (therapist_id, public)
VALUES
  (:therapist_id, :public)
RETURNING id`
