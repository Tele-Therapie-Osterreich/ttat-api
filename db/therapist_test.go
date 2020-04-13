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

func TestLoginTherapist(t *testing.T) {
	RunWithSchema(t, func(pg *PGClient, t *testing.T) {
		LoadDefaultFixture(pg, t)

		var tests = []struct {
			email  string
			userID int
			new    bool
		}{
			{"test1@example.com", 1, false},
			{"newuser@example.com", -1, true},
		}
		for _, test := range tests {
			// TODO: DEAL WITH IMAGES HERE
			th, _, newTh, err := pg.Login(test.email)
			assert.Nil(t, err)
			assert.Equal(t, newTh, test.new, "new therapist mismatch")
			if test.email != "" {
				assert.Equal(t, th.Email, test.email, "email mismatch")
			}
		}
	})
}

// func TestUpdateTherapist(t *testing.T) {
// 	RunWithSchema(t, func(pg *PGClient, t *testing.T) {
// 		loadDefaultFixture(pg, t)

// 		// Check that updating name, etc. work.
// 		therapist, err := pg.TherapistByID("usr_TESTUSER1")
// 		assert.Nil(t, err)
// 		newName := "Update test"
// 		therapist.Name = &newName
// 		assert.Nil(t, pg.UpdateTherapist(therapist))

// 		// Check that updating with bad ID fails.
// 		therapist.ID = "usr_UNKNOWN"
// 		newName = "Update test 2"
// 		therapist.Name = &newName
// 		assert.Equal(t, pg.UpdateTherapist(therapist), ErrTherapistNotFound)

// 		// Check that updating email fails.
// 		therapist, err = pg.TherapistByID("usr_TESTUSER1")
// 		assert.Nil(t, err)
// 		therapist.Email = "new@somewhere.com"
// 		assert.Equal(t, pg.UpdateTherapist(therapist), ErrReadOnlyField,
// 			"update of read-only field!")
// 	})
// }

// func TestDeleteTherapist(t *testing.T) {
// 	RunWithSchema(t, func(pg *PGClient, t *testing.T) {
// 		loadDefaultFixture(pg, t)

// 		var tests = []struct {
// 			id  string
// 			err error
// 		}{
// 			{"usr_TESTUSER1", nil},
// 			{"usr_UNKNOWN", ErrTherapistNotFound},
// 		}
// 		for _, test := range tests {
// 			err := pg.DeleteTherapist(test.id)
// 			assert.Equal(t, err, test.err)
// 			if err != nil {
// 				continue
// 			}
// 			_, err = pg.TherapistByID(test.id)
// 			assert.Equal(t, err, ErrTherapistNotFound, "therapist not deleted")
// 		}
// 	})
// }

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
