package db

import (
	"database/sql"

	"github.com/Tele-Therapie-Osterreich/ttat-api/model"
	"github.com/Tele-Therapie-Osterreich/ttat-api/model/types"
	"github.com/jmoiron/sqlx"
)

// TherapistProfileByTherapistID looks up a therapist's profile by
// their therapist ID.
func (pg *PGClient) TherapistProfileByTherapistID(thID int,
	public bool) (*model.TherapistProfile, error) {
	return therapistProfileByTherapistID(pg.DB, thID, public)
}

func therapistProfileByTherapistID(q sqlx.Queryer, thID int,
	public bool) (*model.TherapistProfile, error) {
	p := &model.TherapistProfile{}
	err := sqlx.Get(q, p, qTherapistProfileByTherapistID, thID, public)
	if err == sql.ErrNoRows {
		return nil, ErrTherapistNotFound
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

const qTherapistProfileByTherapistID = `
SELECT id, therapist_id, public, type,
       name, street_address, city, postcode, country,
       phone, website, languages, short_profile, full_profile
  FROM profiles
 WHERE therapist_id = $1 AND public = $2`

// UpdateTherapistProfile performs profile updating for a therapist.
// If the therapist already has a pending profile, the update profile
// replaces the pending profile. If the therapist does not have a
// pending profile, a new one is generated from the profile passed in.
func (pg *PGClient) UpdateTherapistProfile(thID int, patch []byte) (*model.ImagePatch, error) {
	tx, err := pg.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	p := &model.TherapistProfile{}
	err = tx.Get(p, qEditedTherapistProfileByTherapistID, thID)
	if err == sql.ErrNoRows {
		return nil, ErrTherapistNotFound
	}
	if err != nil {
		return nil, err
	}
	// TODO: DEAL WITH IMAGE PATCH HERE
	image, err := p.Patch(patch)
	if err != nil {
		return nil, err
	}

	// Do the update.
	_, err = tx.NamedExec(qUpdateTherapistProfile, p)
	if err != nil {
		return nil, err
	}
	return image, nil
}

const qEditedTherapistProfileByTherapistID = `
SELECT id, therapist_id, public, type,
       name, street_address, city, postcode, country,
       phone, website, languages, short_profile, full_profile
  FROM profiles
 WHERE id = $1
ORDER BY public ASC
LIMIT 1`

const qUpdateTherapistProfile = `
INSERT INTO profiles
  (therapist_id, public,
   type, name, street_address, city, postcode, country,
   phone, website, languages, short_profile, full_profile)
VALUES
  (:therapist_id, FALSE,
   :type, :name, :street_address, :city, :postcode, :country,
   :phone, :website, :languages, :short_profile, :full_profile)
ON CONFLICT (therapist_id, public)
DO UPDATE
   SET type=:type, name=:name, street_address=:street_address,
       city=:city, postcode=:postcode, country=:country,
       phone=:phone, website=:website, languages=:languages,
       short_profile=:short_profile, full_profile=:full_profile,
       edited_at = now()
 WHERE therapist_id = :therapist_id AND NOT public`

// AbandonTherapistEdits deletes any pending edits profile associated
// with an active therapist.
func (pg *PGClient) AbandonTherapistEdits(thID int) error {
	tx, err := pg.DB.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	th, err := therapistByID(tx, thID)
	if err == sql.ErrNoRows {
		return ErrTherapistNotFound
	}
	if err != nil {
		return err
	}
	if th.Status != types.Active {
		return nil
	}

	result, err := pg.DB.Exec(qAbandonTherapistEdits, thID)
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

const qAbandonTherapistEdits = `
DELETE FROM profiles WHERE therapist_id = $1 AND NOT public`
