package db

import (
	"database/sql"

	"github.com/Tele-Therapie-Osterreich/ttat-api/model"

	// Import Postgres DB driver.
	_ "github.com/lib/pq"
)

// ImageByID looks up an image by its ID.
func (pg *PGClient) ImageByID(id int) (*model.Image, error) {
	image := &model.Image{}
	err := pg.DB.Get(image, imageByID, id)
	if err == sql.ErrNoRows {
		return nil, ErrImageNotFound
	}
	if err != nil {
		return nil, err
	}
	return image, nil
}

const imageByID = `
SELECT id, user_id, extension, data
  FROM images
 WHERE id = $1`

// ImageByUserID looks up a user's image by their user ID.
func (pg *PGClient) ImageByUserID(id int) (*model.Image, error) {
	image := &model.Image{}
	err := pg.DB.Get(image, imageByUserID, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return image, nil
}

const imageByUserID = `
SELECT id, user_id, extension, data
  FROM images
 WHERE user_id = $1`

// UpsertImage inserts or updates a user image.
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

	_, err = pg.UserByID(image.UserID)
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
INSERT INTO images (user_id, extension, data)
             VALUES (:user_id, :extension, :data)
RETURNING id`

const deleteImage = `DELETE FROM images WHERE id = $1`
