package model

import (
	"encoding/json"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/Tele-Therapie-Osterreich/ttat-api/model/types"
)

// Therapist is the database model for therapists who are registered
// with the site.
type Therapist struct {
	// Unique ID of the therapist.
	ID int `db:"id"`

	// Email address used by therapist to log in to account.
	Email string `db:"email"`

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

	// Languages spoken (list of ISO-639 2-letter codes)
	Languages pq.StringArray `db:"languages"`

	// Short profile text.
	ShortProfile *string `db:"short_profile"`

	// Full profile text.
	FullProfile *string `db:"full_profile"`

	// Are the therapist account details approved, pending edits, etc.?
	Status types.ApprovalState `db:"status"`

	// Creation timestamp.
	CreatedAt time.Time `db:"created_at"`
}

// Patch applies a patch represented as a JSON object to a therapist
// value.
func (u *Therapist) Patch(patch []byte) (*ImagePatch, error) {
	updates := map[string]interface{}{}
	err := json.Unmarshal(patch, &updates)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshaling patch")
	}
	roFields := map[string]string{
		"id":         "ID",
		"email":      "email",
		"status":     "approval status",
		"created_at": "creation date",
	}
	for fld, label := range roFields {
		if _, ok := updates[fld]; ok {
			return nil, errors.New("can't patch therapist " + label)
		}
	}
	if err := optStringUpdate(updates, "name", &u.Name); err != nil {
		return nil, err
	}
	if err := optStringUpdate(updates, "street_address", &u.StreetAddress); err != nil {
		return nil, err
	}
	if err := optStringUpdate(updates, "city", &u.City); err != nil {
		return nil, err
	}
	if err := optStringUpdate(updates, "postcode", &u.Postcode); err != nil {
		return nil, err
	}
	if err := optStringUpdate(updates, "country", &u.Country); err != nil {
		return nil, err
	}
	if err := optStringUpdate(updates, "phone", &u.Phone); err != nil {
		return nil, err
	}
	if err := optStringUpdate(updates, "short_profile", &u.ShortProfile); err != nil {
		return nil, err
	}
	if err := optStringUpdate(updates, "full_profile", &u.FullProfile); err != nil {
		return nil, err
	}

	// TODO: DEAL WITH LANGUAGES
	// TODO: DEAL WITH SUB-SPECIALITY LIST.

	return DecodeImagePatch(updates, "photo")
}
