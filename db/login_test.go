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
			name  string
		}{
			{"test1@example.com", false, "User Test1"},
			{"newuser@example.com", true, ""},
		}
		for _, test := range tests {
			// TODO: DEAL WITH IMAGES HERE
			fmt.Println("TestLoginTherapist: email =", test.email,
				"  new =", test.new, "  name =", test.name)
			info, newTh, err := pg.Login(test.email)
			assert.Nil(t, err)
			assert.Equal(t, newTh, test.new, "new therapist mismatch")
			if test.email != "" {
				assert.Equal(t, info.Base.Email, test.email, "email mismatch")
			}
			if test.name != "" {
				assert.Equal(t, *info.Profile.Name, test.name, "name mismatch")
			}
			assert.Equal(t, info.Profile.Public, !test.new)
		}
	})
}
