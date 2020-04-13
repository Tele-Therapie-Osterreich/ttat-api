// +build db

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	InitTestDB()
}

func TestTherapistRetrieval(t *testing.T) {
	RunWithSchema(t, func(pg *PGClient, t *testing.T) {
		LoadDefaultFixture(pg, t)

		var tests = []struct {
			id    int
			err   error
			email string
		}{
			{1, nil, "test1@example.com"},
			{123, ErrTherapistNotFound, ""},
		}
		for _, test := range tests {
			th, err := pg.TherapistByID(test.id)
			assert.Equal(t, err, test.err)
			if err != nil {
				continue
			}
			if assert.NotNil(t, th) {
				assert.Equal(t, th.Email, test.email, "wrong user!")
			}
		}
	})
}

func TestUpdateTherapistSimple(t *testing.T) {
	RunWithSchema(t, func(pg *PGClient, t *testing.T) {
		LoadDefaultFixture(pg, t)

		// Check that updating name, etc. work.
		th, err := pg.TherapistByID(1)
		assert.Nil(t, err)
		newName := "Update test"
		th.Name = &newName
		assert.Nil(t, pg.UpdateTherapist(th))
		thCheck, err := pg.TherapistByID(1)
		assert.Equal(t, *thCheck.Name, newName)

		// Check that updating with bad ID fails.
		th.ID = 123
		newName = "Update test 2"
		th.Name = &newName
		assert.Equal(t, pg.UpdateTherapist(th), ErrTherapistNotFound)

		// Check that updating email fails.
		th, err = pg.TherapistByID(1)
		assert.Nil(t, err)
		th.Email = "new@somewhere.com"
		assert.Equal(t, pg.UpdateTherapist(th), ErrReadOnlyField,
			"update of read-only field!")
	})
}

func TestDeleteTherapist(t *testing.T) {
	RunWithSchema(t, func(pg *PGClient, t *testing.T) {
		LoadDefaultFixture(pg, t)

		var tests = []struct {
			id  int
			err error
		}{
			{1, nil},
			{123, ErrTherapistNotFound},
		}
		for _, test := range tests {
			err := pg.DeleteTherapist(test.id)
			assert.Equal(t, err, test.err)
			if err != nil {
				continue
			}
			_, err = pg.TherapistByID(test.id)
			assert.Equal(t, err, ErrTherapistNotFound, "therapist not deleted")
		}
	})
}

// func TestListUsers(t *testing.T) {
// 	RunWithSchema(t, func(pg *PGClient, t *testing.T) {
// 		loadDefaultFixture(pg, t)

// 		// Get all users: test for count and ordering.
// 		users, err := pg.Users("", 1, 30)
// 		assert.Nil(t, err)
// 		assert.Len(t, users, 10)
// 		assertOrdered(t, users)

// 		// Pagination: test for count and ordering.
// 		users, err = pg.Users("", 1, 4)
// 		assert.Nil(t, err)
// 		assert.Len(t, users, 4)
// 		assertOrdered(t, users)
// 		users, err = pg.Users("", 2, 4)
// 		assert.Nil(t, err)
// 		assert.Len(t, users, 4)
// 		assertOrdered(t, users)

// 		// Search.
// 		users, err = pg.Users("Target", 1, 30)
// 		assert.Nil(t, err)
// 		assert.Len(t, users, 2)
// 		assertOrdered(t, users)
// 	})
// }

// func assertOrdered(t *testing.T, users []model.User) {
// 	for i, u := range users {
// 		if i > 0 {
// 			assert.True(t, u.LastLogin.Before(users[i-1].LastLogin))
// 		}
// 	}
// }
