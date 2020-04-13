package db

import (
	"database/sql"

	"github.com/Tele-Therapie-Osterreich/ttat-api/model"

	// Import Postgres DB driver.
	_ "github.com/lib/pq"
)

// ImageByID looks up an image by its ID.
func (pg *PGClient) ImageByID(imgID int) (*model.Image, error) {
	image := &model.Image{}
	err := pg.DB.Get(image, imageByID, imgID)
	if err == sql.ErrNoRows {
		return nil, ErrImageNotFound
	}
	if err != nil {
		return nil, err
	}
	return image, nil
}

const imageByID = `
SELECT id, therapist_id, extension, data
  FROM images
 WHERE id = $1`

// ImageByTherapistID looks up a therapist's image by their therapist ID.
func (pg *PGClient) ImageByTherapistID(thID int) (*model.Image, error) {
	image := &model.Image{}
	err := pg.DB.Get(image, imageByTherapistID, thID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return image, nil
}

const imageByTherapistID = `
SELECT id, therapist_id, extension, data
  FROM images
 WHERE therapist_id = $1`

// UpsertImage inserts or updates a therapist image.
func (pg *PGClient) UpsertImage(image *model.Image) (*model.Image, error) {
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

	_, err = pg.TherapistByID(image.TherapistID)
	if err != nil {
		return nil, err
	}

	oldID := image.ID
	if oldID != 0 {
		result, err := tx.Exec(deleteImage, oldID)
		if err != nil {
			return nil, err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return nil, err
		}
		if rows != 1 {
			return nil, ErrImageNotFound
		}
	}

	rows, err := tx.NamedQuery(insertImage, image)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	err = rows.Scan(&image.ID)
	if err != nil {
		return nil, err
	}

	return image, nil
}

const insertImage = `
INSERT INTO images (therapist_id, extension, data)
            VALUES (:therapist_id, :extension, :data)
RETURNING id`

const deleteImage = `DELETE FROM images WHERE id = $1`
