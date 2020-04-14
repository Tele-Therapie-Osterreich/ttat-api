// +build db

package db

import (
	"fmt"
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
			email string
			new   bool
		}{
			{"test1@example.com", false},
			{"newuser@example.com", true},
		}
		for _, test := range tests {
			// TODO: DEAL WITH IMAGES HERE
			fmt.Println("TestLoginTherapist: email =", test.email, "  new =", test.new)
			th, _, newTh, err := pg.Login(test.email)
			assert.Nil(t, err)
			assert.Equal(t, newTh, test.new, "new therapist mismatch")
			if test.email != "" {
				assert.Equal(t, th.Email, test.email, "email mismatch")
			}
		}
	})
}
