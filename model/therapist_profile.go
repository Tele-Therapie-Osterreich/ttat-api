package model

import (
	"encoding/json"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/Tele-Therapie-Osterreich/ttat-api/model/types"
)

// TherapistProfile is the database model for profiles of therapists,
// both public and for pending edits.
type TherapistProfile struct {
	// Unique ID of the profile.
	ID int `db:"id"`

	// ID of the therapist this is profile for.
	TherapistID int `db:"therapist_id"`

	// Is this a public profile (true), or a profile for pending edits
	// (false)?
	Public bool `db:"public"`

	// Type of therapist (OT, physio or speech therapist).
	Type types.TherapistType `db:"type"`

	// Name of therapist (optional).
	Name *string `db:"name"`

	// Street address of therapist (optional).
	StreetAddress *string `db:"street_address"`

	// City of therapist's address (optional).
	City *string `db:"city"`

	// Postcode of therapist's address (optional).
	Postcode *string `db:"postcode"`

	// Country of therapist's address (optional).
	Country *string `db:"country"`

	// Contact telephone number for therapist (optional).
	Phone *string `db:"phone"`

	// Therapist website (optional).
	Website *string `db:"website"`

	// Languages spoken (list of ISO-639 2-letter codes)
	Languages pq.StringArray `db:"languages"`

	// Short profile text.
	ShortProfile *string `db:"short_profile"`

	// Full profile text.
	FullProfile *string `db:"full_profile"`

	// Edit timestamp.
	EditedAt time.Time `db:"edited_at"`
}

// Patch applies a patch represented as a JSON object to a therapist
// value.
func (p *TherapistProfile) Patch(patch []byte) (*ImagePatch, error) {
	updates := map[string]interface{}{}
	err := json.Unmarshal(patch, &updates)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshaling patch")
	}
	roFields := map[string]string{
		"id":        "ID",
		"public":    "public flag",
		"edited_at": "creation date",
	}
	for fld, label := range roFields {
		if _, ok := updates[fld]; ok {
			return nil, errors.New("can't patch therapist " + label)
		}
	}
	// TODO: DEAL WITH THERAPIST TYPE
	if err := optStringUpdate(updates, "name", &p.Name); err != nil {
		return nil, err
	}
	if err := optStringUpdate(updates, "street_address", &p.StreetAddress); err != nil {
		return nil, err
	}
	if err := optStringUpdate(updates, "city", &p.City); err != nil {
		return nil, err
	}
	if err := optStringUpdate(updates, "postcode", &p.Postcode); err != nil {
		return nil, err
	}
	if err := optStringUpdate(updates, "country", &p.Country); err != nil {
		return nil, err
	}
	if err := optStringUpdate(updates, "phone", &p.Phone); err != nil {
		return nil, err
	}
	if err := optStringUpdate(updates, "short_profile", &p.ShortProfile); err != nil {
		return nil, err
	}
	if err := optStringUpdate(updates, "full_profile", &p.FullProfile); err != nil {
		return nil, err
	}

	// TODO: DEAL WITH LANGUAGES
	// TODO: DEAL WITH SUB-SPECIALITY LIST.

	// This is an edited profile now.
	p.Public = false

	return DecodeImagePatch(updates, "photo")
}
