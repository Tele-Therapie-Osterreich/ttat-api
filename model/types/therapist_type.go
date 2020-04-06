package types

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
)

// TherapistType is an enumerated type representing the different
// types of therapist: OT, physiotherapist and speech therapist.
type TherapistType uint

const (
	// OccupationalTherapist represents an occupational therapist.
	OccupationalTherapist = iota

	// Physiotherapist represents a physiotherapist.
	Physiotherapist

	// SpeechTherapist represents a speech therapist.
	SpeechTherapist
)

// String converts an approval state to its string representation.
func (t TherapistType) String() string {
	switch t {
	case OccupationalTherapist:
		return "ot"
	case Physiotherapist:
		return "physio"
	case SpeechTherapist:
		return "speech"
	default:
		return "<unknown therapist type>"
	}
}

// FromString does checked conversion from a string to an
// TherapistType.
func (t *TherapistType) FromString(s string) error {
	switch s {
	case "ot":
		*t = OccupationalTherapist
	case "physio":
		*t = Physiotherapist
	case "speech":
		*t = SpeechTherapist
	default:
		return errors.New("unknown therapist type '" + s + "'")
	}
	return nil
}

// MarshalJSON converts an internal approval status to JSON.
func (t TherapistType) MarshalJSON() ([]byte, error) {
	s := t.String()
	if s == "<unknown therapist type>" {
		return nil, errors.New("unknown therapist type")
	}
	return json.Marshal(s)
}

// UnmarshalJSON unmarshals an approval state from a JSON string.
func (t *TherapistType) UnmarshalJSON(d []byte) error {
	var s string
	if err := json.Unmarshal(d, &s); err != nil {
		return errors.Wrap(err, "can't unmarshal therapist type")
	}
	return t.FromString(s)
}

// Scan implements the sql.Scanner interface.
func (t *TherapistType) Scan(src interface{}) error {
	var s string
	switch src := src.(type) {
	case string:
		s = src
	case []byte:
		s = string(src)
	default:
		return errors.New("incompatible type for TherapistType")
	}
	return t.FromString(s)
}

// Value implements the driver.Value interface.
func (t TherapistType) Value() (driver.Value, error) {
	return t.String(), nil
}
