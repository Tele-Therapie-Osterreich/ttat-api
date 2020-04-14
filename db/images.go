package db

import (
	"database/sql"

	"github.com/Tele-Therapie-Osterreich/ttat-api/model"
	"github.com/jmoiron/sqlx"

	// Import Postgres DB driver.
	_ "github.com/lib/pq"
)

// ImageByID looks up an image by its ID.
func (pg *PGClient) ImageByID(imgID int) (*model.Image, error) {
	image := &model.Image{}
	err := pg.DB.Get(image, qImageByID, imgID)
	if err == sql.ErrNoRows {
		return nil, ErrImageNotFound
	}
	if err != nil {
		return nil, err
	}
	return image, nil
}

const qImageByID = `
SELECT id, therapist_id, extension, data
  FROM images
 WHERE id = $1`

// ImageByProfileID looks up a therapist's image by their therapist ID.
func (pg *PGClient) ImageByProfileID(prID int) (*model.Image, error) {
	return imageByProfileID(pg.DB, prID)
}

func imageByProfileID(q sqlx.Queryer, prID int) (*model.Image, error) {
	image := &model.Image{}
	err := sqlx.Get(q, image, qImageByProfileID, prID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return image, nil
}

const qImageByProfileID = `
SELECT id, profile_id, extension, data
  FROM images
 WHERE profile_id = $1`

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

	// TODO: THIS IS WRONG FOR WIPPY REASONS. FIX IT!
	// _, err = pg.TherapistByID(image.ProfileID)
	// if err != nil {
	// 	return nil, err
	// }

	// oldID := image.ID
	// if oldID != 0 {
	// 	result, err := tx.Exec(qDeleteImage, oldID)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	rows, err := result.RowsAffected()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if rows != 1 {
	// 		return nil, ErrImageNotFound
	// 	}
	// }

	// rows, err := tx.NamedQuery(qInsertImage, image)
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()
	// if !rows.Next() {
	// 	return nil, sql.ErrNoRows
	// }

	// err = rows.Scan(&image.ID)
	// if err != nil {
	// 	return nil, err
	// }

	// return image, nil
	return nil, nil
}

const qInsertImage = `
INSERT INTO images (profile_id, extension, data)
            VALUES (:profile_id, :extension, :data)
RETURNING id`

const qDeleteImage = `DELETE FROM images WHERE id = $1`
