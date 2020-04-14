package db

import (
	"database/sql"

	"github.com/Tele-Therapie-Osterreich/ttat-api/model"
)

// PendingEditsByTherapistID looks up the pending edits for a
// therapist by the therapist ID.
func (pg *PGClient) PendingEditsByTherapistID(thID int) (*model.PendingEdits, error) {
	edits := &model.PendingEdits{}
	err := pg.DB.Get(edits, pendingEditsByTherapistID, thID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return edits, nil
}

const pendingEditsByTherapistID = `
SELECT id, therapist_id, patch, edited_at
  FROM therapist_pending_edits
 WHERE therapist_id = $1`

// AddPendingEdits adds some pending edits for a therapist, merging
// with any existing pending edits for the therapist, and updating the
// therapist's approval state.
func (pg *PGClient) AddPendingEdits(thID int, patch []byte) (*model.PendingEdits, error) {
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

	edits := &model.PendingEdits{}
	err = tx.Get(edits, addPendingEdits, thID, patch)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(setTherapistPending, thID)
	if err != nil {
		return nil, err
	}

	return edits, nil
}

const addPendingEdits = `
INSERT INTO therapist_pending_edits AS e (therapist_id, patch) VALUES ($1, $2)
  ON CONFLICT (therapist_id) DO UPDATE
    SET patch = COALESCE(e.patch || EXCLUDED.patch, EXCLUDED.patch)
    WHERE e.therapist_id = $1
RETURNING id, therapist_id, patch, edited_at
`

const setTherapistPending = `
UPDATE therapists SET status = 'edits_pending'
 WHERE id = $1 AND status <> 'new'`
