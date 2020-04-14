package model

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

// PendingEdits is the data model representing edits to a therapist's
// profile that are not yet publicly visible.
type PendingEdits struct {
	// Unique ID of the therapist.
	ID int `db:"id"`

	// ID of therapist this set of edits is for.
	TherapistID int `db:"therapist_id"`

	// JSON representation of the edits.
	Patch types.JSONText `db:"patch"`

	// Edit timestamp.
	EditedAt time.Time `db:"edited_at"`
}
