package types

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
)

// ApprovalState is an enumerated type representing the current status
// of a therapist account. Accounts can be: new, approved, edits
// pending or suspended.
type ApprovalState uint

const (
	// New represents a therapist user account after its initial
	// creation: it is not yet generally visible, and must be approved
	// by an administrator.
	New = iota

	// Approved represents a therapist account that has been approved
	// for release by an administrator, and is now visible to all users.
	Approved

	// EditsPending represents a therapist account that is otherwise
	// approved, but that has edits pending that need to be approved
	// before they are released to general visibility.
	EditsPending

	// Suspended represents a therapist account that an administrator
	// suspended for some reason, so that it is not visible to any users
	// apart from its owner (for editing) and administrators.
	Suspended
)

// String converts an approval state to its string representation.
func (a ApprovalState) String() string {
	switch a {
	case New:
		return "new"
	case Approved:
		return "approved"
	case EditsPending:
		return "edits_pending"
	case Suspended:
		return "suspended"
	default:
		return "<unknown approval state>"
	}
}

// FromString does checked conversion from a string to an
// ApprovalState.
func (a *ApprovalState) FromString(s string) error {
	switch s {
	case "new":
		*a = New
	case "approved":
		*a = Approved
	case "edits_pending":
		*a = EditsPending
	case "suspended":
		*a = Suspended
	default:
		return errors.New("unknown approval state '" + s + "'")
	}
	return nil
}

// MarshalJSON converts an internal approval status to JSON.
func (a ApprovalState) MarshalJSON() ([]byte, error) {
	s := a.String()
	if s == "<unknown approval state>" {
		return nil, errors.New("unknown approval state")
	}
	return json.Marshal(s)
}

// UnmarshalJSON unmarshals an approval state from a JSON string.
func (a *ApprovalState) UnmarshalJSON(d []byte) error {
	var s string
	if err := json.Unmarshal(d, &s); err != nil {
		return errors.Wrap(err, "can't unmarshal approval state")
	}
	return a.FromString(s)
}

// Scan implements the sql.Scanner interface.
func (a *ApprovalState) Scan(src interface{}) error {
	var s string
	switch src := src.(type) {
	case string:
		s = src
	case []byte:
		s = string(src)
	default:
		return errors.New("incompatible type for ApprovalState")
	}
	return a.FromString(s)
}

// Value implements the driver.Value interface.
func (a ApprovalState) Value() (driver.Value, error) {
	return a.String(), nil
}
