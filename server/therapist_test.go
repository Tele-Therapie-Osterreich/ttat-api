// +build !db

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tele-Therapie-Osterreich/ttat-api/db"
	"github.com/Tele-Therapie-Osterreich/ttat-api/mocks"
	"github.com/Tele-Therapie-Osterreich/ttat-api/model"
	"github.com/Tele-Therapie-Osterreich/ttat-api/model/types"
	"github.com/stretchr/testify/assert"
)

func TestTherapistDetailUnknownTherapistID(t *testing.T) {
	d, m, r := mockServer()
	d.On("TherapistInfoByID", TestTherapistID, db.PublicOnly).Return(nil, db.ErrTherapistNotFound)

	rr := apiTest(t, r, "GET", testTherapistURL("/therapist/%d"), nil)

	assert.Equal(t, rr.Code, http.StatusNotFound)
	d.AssertExpectations(t)
	m.AssertExpectations(t)
}

func TestTherapistDetailNoImage(t *testing.T) {
	d, m, r := mockServer()
	setupTherapist(d, false, db.PublicOnly)

	rr := apiTest(t, r, "GET", testTherapistURL("/therapist/%d"), nil)

	assert.Equal(t, rr.Code, http.StatusOK)
	checkView(t, rr, nil)
	d.AssertExpectations(t)
	m.AssertExpectations(t)
}

func TestTherapistDetailWithImage(t *testing.T) {
	d, m, r := mockServer()
	setupTherapist(d, true, db.PublicOnly)

	rr := apiTest(t, r, "GET", testTherapistURL("/therapist/%d"), nil)

	assert.Equal(t, rr.Code, http.StatusOK)
	checkView(t, rr, &viewOptions{photo: testImageURL()})
	d.AssertExpectations(t)
	m.AssertExpectations(t)
}

func TestSelfDetail(t *testing.T) {
	tests := []struct {
		auth bool
		code int
	}{
		{false, http.StatusNotFound},
		{true, http.StatusOK},
	}
	for _, test := range tests {
		fmt.Println("TestSelfDetail: auth =", test.auth, "  code =", test.code)
		d, m, r := mockServer()
		if test.auth {
			setupSession(d)
			setupTherapist(d, test.auth, db.PreferPublic)
		}

		session := ""
		if test.auth {
			session = TestSession
		}
		rr := apiTest(t, r, "GET", "/me", &apiOptions{session: session})

		assert.Equal(t, rr.Code, test.code)
		if test.auth {
			checkView(t, rr, &viewOptions{photo: testImageURL()})
		}
		d.AssertExpectations(t)
		m.AssertExpectations(t)
	}
}

func TestSelfDetailPending(t *testing.T) {
	d, m, r := mockServer()
	setupSession(d)
	setupTherapist(d, true, db.PendingOnly)

	rr := apiTest(t, r, "GET", "/me/pending", &apiOptions{session: TestSession})

	assert.Equal(t, rr.Code, http.StatusOK)
	checkView(t, rr, &viewOptions{photo: testImageURL(), edited: true})
	d.AssertExpectations(t)
	m.AssertExpectations(t)
}

func TestSelfDelete(t *testing.T) {
	tests := []struct {
		auth    bool
		deleted bool
		code    int
	}{
		{false, false, http.StatusNotFound},
		{true, true, http.StatusNoContent},
	}
	for _, test := range tests {
		d, m, r := mockServer()
		if test.auth {
			setupSession(d)
			setupDelete(d)
		}

		session := ""
		if test.auth {
			session = TestSession
		}
		rr := apiTest(t, r, "DELETE", "/me", &apiOptions{session: session})

		assert.Equal(t, rr.Code, test.code)
		d.AssertExpectations(t)
		m.AssertExpectations(t)
	}
}

// func TestSelfUpdateSimple(t *testing.T) {
// 	tests := []struct {
// 		auth    bool
// 		updated bool
// 		json    string
// 		code    int
// 	}{
// 		{false, false, `{"name": "Name Changed"}`, http.StatusNotFound},
// 		{true, false, `{"bad-json": `, http.StatusBadRequest},
// 		{true, false, `{"email": "new@somewhere.com"}`, http.StatusBadRequest},
// 		{true, true, `{"name": "Name Changed"}`, http.StatusOK},
// 	}
// 	for _, test := range tests {
// 		d, m, r := mockServer()
// 		therapistSetup(d, &setupOptions{
//      if test.auth {
//        setupSession(d)
//      }
// 			needsAuth: true,
// 			update:    test.updated,
// 		})

// 		session := ""
// 		if test.auth {
// 			session = TestSession
// 		}
// 		rr := apiTest(t, r, "PATCH", "/me",
// 			&apiOptions{
// 				session:  session,
// 				bodyJSON: []byte(test.json),
// 			})

// 		assert.Equal(t, rr.Code, test.code)
// 		d.AssertExpectations(t)
// 		m.AssertExpectations(t)
// 	}
// }

func setupSession(d *mocks.DB) {
	id := TestTherapistID
	d.On("LookupSession", TestSession).Return(&id, nil)
}

func setupTherapist(d *mocks.DB, image bool, profile db.ProfileSelection) {
	th := &model.Therapist{
		ID:     TestTherapistID,
		Email:  TestEmail,
		Status: types.Active,
	}
	name := TestName
	public := true
	if profile == db.PendingOnly {
		public = false
		name = "Name Changed"
	}
	p := &model.TherapistProfile{
		ID:          TestProfileID,
		TherapistID: TestTherapistID,
		Public:      public,
		Name:        &name,
		Type:        types.OccupationalTherapist,
	}
	var i *model.Image
	if image {
		i = &model.Image{
			ID:        TestImageID,
			ProfileID: TestProfileID,
			Extension: TestImageExtension,
			Data:      []byte{1, 2, 3, 4},
		}
	}
	info := &model.TherapistInfo{
		Base:              th,
		Profile:           p,
		Image:             i,
		HasPublicProfile:  true,
		HasPendingProfile: false,
	}
	d.On("TherapistInfoByID", TestTherapistID, profile).Return(info, nil)
}

func setupDelete(d *mocks.DB) {
	d.On("DeleteTherapist", TestTherapistID).Return(nil)
}

// func therapistSetup(d *mocks.DB, opts *setupOptions) {
// 	if opts != nil && opts.needsAuth && !opts.session {
// 		return
// 	}

// 	if opts == nil || !opts.delete {
// 		d.On("ImageByTherapistID", TestTherapistID).Return(nil, nil)
// 	}

// 	if opts != nil && opts.pendingEdits {

// 	}

// 	if opts != nil && opts.update {
// 		d.On("UpdateTherapist", mock.AnythingOfType("*model.Therapist")).Return(nil)
// 	}
// }

type viewOptions struct {
	photo  string
	edited bool
}

func checkView(t *testing.T, rr *httptest.ResponseRecorder, opts *viewOptions) {
	t.Helper()
	v := model.TherapistFullView{}
	err := json.Unmarshal(rr.Body.Bytes(), &v)
	assert.Nil(t, err)
	assert.Equal(t, v.ID, TestTherapistID)
	assert.Equal(t, v.Email, TestEmail)
	assert.NotNil(t, v.Name)
	if opts == nil || !opts.edited {
		assert.Equal(t, *v.Name, TestName)
	} else {
		assert.Equal(t, *v.Name, "Name Changed")
	}
	if opts == nil || opts.photo == "" {
		assert.Nil(t, v.Photo)
	} else {
		assert.NotNil(t, v.Photo)
		assert.Equal(t, *v.Photo, opts.photo)
	}
}
