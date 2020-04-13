// +build !db

package server

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Tele-Therapie-Osterreich/ttat-api/db"
	"github.com/Tele-Therapie-Osterreich/ttat-api/messages"
	"github.com/Tele-Therapie-Osterreich/ttat-api/model"
)

func TestRequestLoginEmailBadJSON(t *testing.T) {
	d, m, r := mockServer()

	rr := apiTest(t, r, "POST", "/auth/request-login-email",
		&apiOptions{bodyJSON: []byte(`{"broken":}`)})

	assert.Equal(t, rr.Code, http.StatusBadRequest)
	resp := messages.APIError{}
	assert.Nil(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, resp.Code, messages.InvalidJSONBody)

	d.AssertExpectations(t)
	m.AssertExpectations(t)
}

func TestRequestLoginEmailBadEmail(t *testing.T) {
	d, m, r := mockServer()

	rr := apiTest(t, r, "POST", "/auth/request-login-email",
		&apiOptions{body: &messages.ReqLoginEmailRequest{
			Email:    "bob@x",
			Language: "en",
		}})

	assert.Equal(t, rr.Code, http.StatusBadRequest)
	resp := messages.APIError{}
	assert.Nil(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, resp.Code, messages.InvalidEmailAddress)

	d.AssertExpectations(t)
	m.AssertExpectations(t)
}

func TestRequestLoginEmail(t *testing.T) {
	d, m, r := mockServer()

	d.On("CreateLoginToken", TestEmail, "en").Return(TestToken, nil)
	m.On("Send", "login-email-request", TestEmail, "en",
		map[string]string{"login_token": TestToken}).Return(nil)

	rr := apiTest(t, r, "POST", "/auth/request-login-email",
		&apiOptions{body: &messages.ReqLoginEmailRequest{
			Email:    TestEmail,
			Language: "en",
		}})

	assert.Equal(t, rr.Code, http.StatusNoContent)
	assert.Empty(t, rr.Body.String())

	d.AssertExpectations(t)
	m.AssertExpectations(t)
}

func TestLoginBadLoginToken(t *testing.T) {
	d, m, r := mockServer()

	d.On("CheckLoginToken", "bad-token").Return("", "", db.ErrLoginTokenNotFound)

	rr := apiTest(t, r, "POST", "/auth/login",
		&apiOptions{body: &messages.LoginRequest{
			LoginToken: "bad-token",
		}})

	assert.Equal(t, rr.Code, http.StatusBadRequest)
	resp := messages.APIError{}
	assert.Nil(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, resp.Code, messages.UnknownLoginToken)

	d.AssertExpectations(t)
	m.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
	d, m, r := mockServer()

	d.On("CheckLoginToken", TestToken).Return(TestEmail, "en", nil)
	u := model.Therapist{
		ID:    TestTherapistID,
		Email: TestEmail,
	}
	d.On("Login", TestEmail).Return(&u, nil, true, nil)
	d.On("CreateSession", TestTherapistID).Return(TestSession, nil)

	rr := apiTest(t, r, "POST", "/auth/login",
		&apiOptions{body: &messages.LoginRequest{
			LoginToken: TestToken,
		}})
	assert.Equal(t, rr.Code, http.StatusOK)

	resp := messages.LoginResponse{}
	assert.Nil(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, resp.Profile.ID, TestTherapistID)

	d.AssertExpectations(t)
	m.AssertExpectations(t)
}
