package db

import (
	"database/sql"

	"github.com/Tele-Therapie-Osterreich/ttat-api/model"
	"github.com/jmoiron/sqlx"
)

// TherapistInfoByEmail returns the therapist model, public or pending
// profile and image for a given therapist ID.
func (pg *PGClient) TherapistInfoByEmail(email string) (*model.TherapistInfo, error) {
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

	return therapistInfoByEmail(tx, email)
}

func therapistInfoByEmail(q sqlx.Queryer, email string) (*model.TherapistInfo, error) {
	th, err := therapistByEmail(q, email)
	if err != nil {
		return nil, err
	}

	return collectTherapistInfo(q, th, PreferPublic)
}

// TherapistInfoByID returns the therapist model, public or pending
// profile and image for a given therapist ID.
func (pg *PGClient) TherapistInfoByID(thID int,
	profile ProfileSelection) (*model.TherapistInfo, error) {
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

	th, err := therapistByID(tx, thID)
	if err != nil {
		return nil, err
	}

	return collectTherapistInfo(tx, th, profile)
}

func selectProfile(public bool) ProfileSelection {
	if public {
		return PublicOnly
	}
	return PendingOnly
}

// TherapistInfoByID returns the therapist model, public or pending
// profile and image for a given therapist ID.
func collectTherapistInfo(q sqlx.Queryer, th *model.Therapist,
	selector ProfileSelection) (*model.TherapistInfo, error) {
	var p *model.TherapistProfile
	var err error
	switch selector {
	case PublicOnly:
		p, err = therapistProfileByTherapistID(q, th.ID, true)
	case PendingOnly:
		p, err = therapistProfileByTherapistID(q, th.ID, false)
	case PreferPublic:
		p, err = therapistProfilePreferPublic(q, th.ID)
	}
	if err != nil {
		return nil, err
	}

	image, err := imageByProfileID(q, p.ID)
	if err != nil {
		return nil, err
	}

	statuses := []bool{}
	err = sqlx.Select(q, &statuses, qProfileTypes, th.ID)
	if err != nil {
		return nil, err
	}
	hasPublic := false
	hasPending := false
	for _, status := range statuses {
		if status {
			hasPublic = true
		} else {
			hasPending = true
		}
	}

	return &model.TherapistInfo{
		Base:              th,
		Profile:           p,
		Image:             image,
		HasPublicProfile:  hasPublic,
		HasPendingProfile: hasPending,
	}, nil
}

const qProfileTypes = `
SELECT public FROM profiles WHERE therapist_id = $1`

func therapistProfilePreferPublic(q sqlx.Queryer, thID int) (*model.TherapistProfile, error) {
	p := &model.TherapistProfile{}
	err := sqlx.Get(q, p, qTherapistProfilePreferPublic, thID)
	if err == sql.ErrNoRows {
		return nil, ErrTherapistNotFound
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

const qTherapistProfilePreferPublic = `
SELECT id, therapist_id, public, type,
       name, street_address, city, postcode, country,
       phone, website, languages, short_profile, full_profile
  FROM profiles
 WHERE therapist_id = $1
ORDER BY public DESC LIMIT 1`
