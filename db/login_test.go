// +build db

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	InitTestDB()
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
