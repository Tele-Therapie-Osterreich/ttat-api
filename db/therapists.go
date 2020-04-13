package db

import (
	"database/sql"
	"fmt"

	"github.com/Tele-Therapie-Osterreich/ttat-api/model"

	// Import Postgres DB driver.
	_ "github.com/lib/pq"
)

// TherapistByID looks up a therapist by their therapist ID.
func (pg *PGClient) TherapistByID(id int) (*model.Therapist, error) {
	th := &model.Therapist{}
	err := pg.DB.Get(th, therapistByID, id)
	if err == sql.ErrNoRows {
		return nil, ErrTherapistNotFound
	}
	if err != nil {
		return nil, err
	}
	return th, nil
}

const therapistByID = `
SELECT id, email, type, name,
       street_address, city, postcode, country,
       phone, languages, short_profile, full_profile,
       status, created_at
  FROM therapists
 WHERE id = $1`

// UpdateTherapist updates the therapist's details in the database.
// The id, email, status and created_at fields are read-only using
// this method.
func (pg *PGClient) UpdateTherapist(th *model.Therapist) error {
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

	check := &model.Therapist{}
	err = tx.Get(check, therapistByID, th.ID)
	if err == sql.ErrNoRows {
		return ErrTherapistNotFound
	}
	if err != nil {
		return err
	}

	// Check read-only fields.
	if th.Email != check.Email || th.Status != check.Status ||
		th.CreatedAt != check.CreatedAt {
		return ErrReadOnlyField
	}

	// Fix other field values: name fields should use empty strings for
	// null values.
	empty := ""
	if th.Name == nil {
		th.Name = &empty
	}
	if th.StreetAddress == nil {
		th.StreetAddress = &empty
	}
	if th.City == nil {
		th.City = &empty
	}
	if th.Postcode == nil {
		th.Postcode = &empty
	}
	if th.Country == nil {
		th.Country = &empty
	}
	if th.Phone == nil {
		th.Phone = &empty
	}
	if th.ShortProfile == nil {
		th.ShortProfile = &empty
	}
	if th.FullProfile == nil {
		th.FullProfile = &empty
	}

	// Do the update.
	_, err = tx.NamedExec(updateTherapist, th)
	if err != nil {
		return err
	}
	return nil
}

const updateTherapist = `
UPDATE therapists
   SET type=:type, name=:name, street_address=:street_address,
       city=:city, postcode=:postcode, country=:country,
       phone=:phone, languages=:languages,
       short_profile=:short_profile, full_profile=:full_profile
 WHERE id = :id`

// DeleteTherapist deletes the given therapist account.
func (pg *PGClient) DeleteTherapist(thID int) error {
	fmt.Println("===> DeleteTherapist")
	result, err := pg.DB.Exec(deleteTherapist, thID)
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

const deleteTherapist = "DELETE FROM therapists WHERE id = $1"
