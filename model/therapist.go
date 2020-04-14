package model

import (
	"time"

	"github.com/Tele-Therapie-Osterreich/ttat-api/model/types"
)

// Therapist is the database model for therapists who are registered
// with the site.
type Therapist struct {
	// Unique ID of the therapist.
	ID int `db:"id"`

	// Email address used by therapist to log in to account.
	Email string `db:"email"`

	// Are the therapist account details approved, pending edits, etc.?
	Status types.ApprovalState `db:"status"`

	// Last login timestamp.
	LastLoginAt time.Time `db:"last_login_at"`

	// Creation timestamp.
	CreatedAt time.Time `db:"created_at"`
}
