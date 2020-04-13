// +build !db

package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tele-Therapie-Osterreich/ttat-api/db"
	"github.com/Tele-Therapie-Osterreich/ttat-api/mocks"
	"github.com/Tele-Therapie-Osterreich/ttat-api/model"
	"github.com/Tele-Therapie-Osterreich/ttat-api/model/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTherapistDetailUnknownTherapistID(t *testing.T) {
	d, m, r := mockServer()
	d.On("TherapistByID", TestTherapistID).Return(nil, db.ErrTherapistNotFound)

	rr := apiTest(t, r, "GET", testTherapistURL("/therapist/%d"), nil)

	assert.Equal(t, rr.Code, http.StatusNotFound)
	d.AssertExpectations(t)
	m.AssertExpectations(t)
}

func TestTherapistDetailNoImage(t *testing.T) {
	d, m, r := mockServer()
	therapistSetup(d, nil)

	rr := apiTest(t, r, "GET", testTherapistURL("/therapist/%d"), nil)

	assert.Equal(t, rr.Code, http.StatusOK)
	checkProfile(t, rr, nil)
	d.AssertExpectations(t)
	m.AssertExpectations(t)
}

func TestTherapistDetailWithImage(t *testing.T) {
	d, m, r := mockServer()
	therapistSetup(d, &setupOptions{image: true})

	rr := apiTest(t, r, "GET", testTherapistURL("/therapist/%d"), nil)

	assert.Equal(t, rr.Code, http.StatusOK)
	checkProfile(t, rr, &profileOptions{photo: testImageURL()})
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
		d, m, r := mockServer()
		therapistSetup(d, &setupOptions{
			session:   test.auth,
			image:     test.auth,
			needsAuth: true,
		})

		session := ""
		if test.auth {
			session = TestSession
		}
		rr := apiTest(t, r, "GET", "/me", &apiOptions{session: session})

		assert.Equal(t, rr.Code, test.code)
		if test.auth {
			checkProfile(t, rr, &profileOptions{photo: testImageURL()})
		}
		d.AssertExpectations(t)
		m.AssertExpectations(t)
	}
}

// func TestSelfDetailPending(t *testing.T) {
// 	d, m, r := mockServer()
// 	therapistSetup(d, &setupOptions{image: true, session: true, pendingEdits: true})

// 	rr := apiTest(t, r, "GET", "/me?status=pending", &apiOptions{session: TestSession})

// 	assert.Equal(t, rr.Code, http.StatusOK)
// 	checkProfile(t, rr, &profileOptions{photo: testImageURL(), edited: true})
// 	d.AssertExpectations(t)
// 	m.AssertExpectations(t)
// }

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
		therapistSetup(d, &setupOptions{
			session:   test.auth,
			needsAuth: true,
			delete:    test.deleted,
		})

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

func TestSelfUpdateSimple(t *testing.T) {
	tests := []struct {
		auth    bool
		updated bool
		json    string
		code    int
	}{
		{false, false, `{"name": "Name Changed"}`, http.StatusNotFound},
		{true, false, `{"bad-json": `, http.StatusBadRequest},
		{true, false, `{"email": "new@somewhere.com"}`, http.StatusBadRequest},
		{true, true, `{"name": "Name Changed"}`, http.StatusOK},
	}
	for _, test := range tests {
		d, m, r := mockServer()
		therapistSetup(d, &setupOptions{
			session:   test.auth,
			needsAuth: true,
			update:    test.updated,
		})

		session := ""
		if test.auth {
			session = TestSession
		}
		rr := apiTest(t, r, "PATCH", "/me",
			&apiOptions{
				session:  session,
				bodyJSON: []byte(test.json),
			})

		assert.Equal(t, rr.Code, test.code)
		d.AssertExpectations(t)
		m.AssertExpectations(t)
	}
}

type setupOptions struct {
	needsAuth    bool
	session      bool
	image        bool
	pendingEdits bool
	delete       bool
	update       bool
}

func therapistSetup(d *mocks.DB, opts *setupOptions) {
	if opts != nil && opts.session {
		id := TestTherapistID
		d.On("LookupSession", TestSession).Return(&id, nil)
	}

	if opts != nil && opts.needsAuth && !opts.session {
		return
	}

	if opts == nil || !opts.delete {
		name := TestName
		u := model.Therapist{
			ID:     TestTherapistID,
			Email:  TestEmail,
			Name:   &name,
			Type:   types.OccupationalTherapist,
			Status: types.Approved,
		}
		d.On("TherapistByID", TestTherapistID).Return(&u, nil)
	}

	if opts != nil && opts.image {
		i := model.Image{
			ID:          TestImageID,
			TherapistID: TestTherapistID,
			Extension:   TestImageExtension,
			Data:        []byte{1, 2, 3, 4},
		}
		d.On("ImageByTherapistID", TestTherapistID).Return(&i, nil)
	}
	if opts == nil || !opts.delete {
		d.On("ImageByTherapistID", TestTherapistID).Return(nil, nil)
	}

	if opts != nil && opts.pendingEdits {

	}

	if opts != nil && opts.delete {
		d.On("DeleteTherapist", TestTherapistID).Return(nil)
	}

	if opts != nil && opts.update {
		d.On("UpdateTherapist", mock.AnythingOfType("*model.Therapist")).Return(nil)
	}
}

type profileOptions struct {
	photo  string
	edited bool
}

func checkProfile(t *testing.T, rr *httptest.ResponseRecorder,
	opts *profileOptions) {
	t.Helper()
	profile := model.TherapistFullProfile{}
	err := json.Unmarshal(rr.Body.Bytes(), &profile)
	assert.Nil(t, err)
	assert.Equal(t, profile.ID, TestTherapistID)
	assert.Equal(t, profile.Email, TestEmail)
	assert.NotNil(t, profile.Name)
	if opts == nil || !opts.edited {
		assert.Equal(t, *profile.Name, TestName)
	} else {
		assert.Equal(t, *profile.Name, "Name Changed")
	}
	if opts == nil || opts.photo == "" {
		assert.Nil(t, profile.Photo)
	} else {
		assert.NotNil(t, profile.Photo)
		assert.Equal(t, *profile.Photo, opts.photo)
	}
}
