// +build db

package db

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	InitTestDB()
}

func TestPendingEditRetrieval(t *testing.T) {
	RunWithSchema(t, func(pg *PGClient, t *testing.T) {
		LoadDefaultFixture(pg, t)

		var tests = []struct {
			thID      int
			nilReturn bool
		}{
			{2, false},
			{123, true},
		}
		for _, test := range tests {
			fmt.Println("TestPendingEditRetrieval: therapist_id =", test.thID)

			edits, err := pg.PendingEditsByTherapistID(test.thID)

			assert.Nil(t, err)
			if test.nilReturn {
				assert.Nil(t, edits)
				continue
			}

			assert.NotNil(t, edits)
			assert.GreaterOrEqual(t, len(edits.Patch), 5)
		}
	})
}

func TestAddPendingEdits(t *testing.T) {
	RunWithSchema(t, func(pg *PGClient, t *testing.T) {
		LoadDefaultFixture(pg, t)

		patch := []byte(`{"name": "Changed Name", "short_profile": "Edited profile"}`)
		edits, err := pg.AddPendingEdits(1, patch)
		assert.Nil(t, err)
		fmt.Println("TestAddPendingEdits (1): edits.Patch =", string(edits.Patch))
		check1 := map[string]interface{}{}
		err = json.Unmarshal(edits.Patch, &check1)
		assert.Nil(t, err)
		assert.Len(t, check1, 2)
		assert.Contains(t, check1, "name")
		assert.Contains(t, check1, "short_profile")
		assert.Equal(t, check1["name"], "Changed Name")
		assert.Equal(t, check1["short_profile"], "Edited profile")

		edits, err = pg.AddPendingEdits(2, patch)
		assert.Nil(t, err)
		fmt.Println("TestAddPendingEdits (2): edits.Patch =", string(edits.Patch))
		check2 := map[string]interface{}{}
		err = json.Unmarshal(edits.Patch, &check2)
		assert.Nil(t, err)
		assert.Len(t, check2, 11)
		assert.Contains(t, check2, "type")
		assert.Contains(t, check2, "name")
		assert.Contains(t, check2, "street_address")
		assert.Contains(t, check2, "city")
		assert.Contains(t, check2, "postcode")
		assert.Contains(t, check2, "country")
		assert.Contains(t, check2, "phone")
		assert.Contains(t, check2, "website")
		assert.Contains(t, check2, "languages")
		assert.Contains(t, check2, "short_profile")
		assert.Contains(t, check2, "full_profile")
		assert.Equal(t, check2["name"], "Changed Name")
		assert.Equal(t, check2["short_profile"], "Edited profile")

		edits, err = pg.AddPendingEdits(3, patch)
		assert.Nil(t, err)
		fmt.Println("TestAddPendingEdits (3): edits.Patch =", string(edits.Patch))
		check3 := map[string]interface{}{}
		err = json.Unmarshal(edits.Patch, &check3)
		assert.Nil(t, err)
		assert.Len(t, check3, 2)
		assert.Contains(t, check3, "name")
		assert.Contains(t, check3, "short_profile")
		assert.Equal(t, check3["name"], "Changed Name")
		assert.Equal(t, check3["short_profile"], "Edited profile")
	})
}
