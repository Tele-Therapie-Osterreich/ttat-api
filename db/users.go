package db

import (
	"database/sql"

	"github.com/Tele-Therapie-Osterreich/ttat-api/model"

	// Import Postgres DB driver.
	_ "github.com/lib/pq"
)

// UserByID looks up a user by their user ID.
func (pg *PGClient) UserByID(id int) (*model.User, error) {
	user := &model.User{}
	err := pg.DB.Get(user, userByID, id)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

const userByID = `
SELECT id, email, type, name,
       street_address, city, postcode, country,
       phone, short_profile, full_profile,
       status, created_at
  FROM users
 WHERE id = $1`

// UpdateUser updates the user's details in the database. The id,
// email, status and created_at fields are read-only using this
// method.
func (pg *PGClient) UpdateUser(user *model.User) error {
	tx, err := pg.DB.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	check := &model.User{}
	err = tx.Get(check, userByID, user.ID)
	if err == sql.ErrNoRows {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}

	// Check read-only fields.
	if user.Email != check.Email || user.Status != check.Status ||
		user.CreatedAt != check.CreatedAt {
		return ErrReadOnlyField
	}

	// Fix other field values: name fields should use empty strings for
	// null values.
	empty := ""
	if user.Name == nil {
		user.Name = &empty
	}
	if user.StreetAddress == nil {
		user.StreetAddress = &empty
	}
	if user.City == nil {
		user.City = &empty
	}
	if user.Postcode == nil {
		user.Postcode = &empty
	}
	if user.Country == nil {
		user.Country = &empty
	}
	if user.Phone == nil {
		user.Phone = &empty
	}
	if user.ShortProfile == nil {
		user.ShortProfile = &empty
	}
	if user.FullProfile == nil {
		user.FullProfile = &empty
	}

	// Do the update.
	_, err = tx.NamedExec(updateUser, user)
	if err != nil {
		return err
	}
	return nil
}

const updateUser = `
UPDATE users SET type=:type, name=:name, street_address=:street_address,
                 city=:city, postcode=:postcode, country=:country,
                 phone=:phone, short_profile=:short_profile,
                 full_profile=:full_profile
 WHERE id = :id`

// DeleteUser deletes the given user account.
func (pg *PGClient) DeleteUser(id int) error {
	result, err := pg.DB.Exec(deleteUser, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return ErrUserNotFound
	}
	return nil
}

const deleteUser = "DELETE FROM users WHERE id = $1"
