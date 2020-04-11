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

const TestEmail = "test@example.com"
const TestToken = "test-token"
const TestUserID = 1
const TestSession = "test-session"

func TestRequestLoginEmailBadJSON(t *testing.T) {
	s, d, m := MockServer()

	rr := apiTest(t, SimpleHandler(s.RequestLoginEmail),
		"POST", []byte(`{"broken":}`))

	assert.Equal(t, rr.Code, http.StatusBadRequest)
	resp := messages.APIError{}
	assert.Nil(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, resp.Code, messages.InvalidJSONBody)

	d.AssertExpectations(t)
	m.AssertExpectations(t)
}

func TestRequestLoginEmailBadEmail(t *testing.T) {
	s, d, m := MockServer()

	reqLoginEmail := messages.ReqLoginEmailRequest{
		Email:    "bob@x",
		Language: "en",
	}
	body, _ := json.Marshal(reqLoginEmail)

	rr := apiTest(t, SimpleHandler(s.RequestLoginEmail), "POST", body)

	assert.Equal(t, rr.Code, http.StatusBadRequest)
	resp := messages.APIError{}
	assert.Nil(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, resp.Code, messages.InvalidEmailAddress)

	d.AssertExpectations(t)
	m.AssertExpectations(t)
}

func TestRequestLoginEmail(t *testing.T) {
	s, d, m := MockServer()

	d.On("CreateLoginToken", TestEmail, "en").Return(TestToken, nil)
	m.On("Send", "login-email-request", TestEmail, "en",
		map[string]string{"login_token": TestToken}).Return(nil)

	reqLoginEmail := messages.ReqLoginEmailRequest{
		Email:    TestEmail,
		Language: "en",
	}
	body, _ := json.Marshal(reqLoginEmail)

	rr := apiTest(t, SimpleHandler(s.RequestLoginEmail), "POST", body)

	assert.Equal(t, rr.Code, http.StatusNoContent)
	assert.Empty(t, rr.Body.String())

	d.AssertExpectations(t)
	m.AssertExpectations(t)
}

func TestLoginBadLoginToken(t *testing.T) {
	s, d, m := MockServer()

	d.On("CheckLoginToken", "bad-token").Return("", "", db.ErrLoginTokenNotFound)

	login := messages.LoginRequest{
		LoginToken: "bad-token",
	}
	body, _ := json.Marshal(login)

	rr := apiTest(t, SimpleHandler(s.Login), "POST", body)

	assert.Equal(t, rr.Code, http.StatusBadRequest)
	resp := messages.APIError{}
	assert.Nil(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, resp.Code, messages.UnknownLoginToken)

	d.AssertExpectations(t)
	m.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
	s, d, m := MockServer()

	d.On("CheckLoginToken", TestToken).Return(TestEmail, "en", nil)
	u := model.User{
		ID:    TestUserID,
		Email: TestEmail,
	}
	d.On("LoginUser", TestEmail).Return(&u, nil, true, nil)
	d.On("CreateSession", TestUserID).Return(TestSession, nil)

	login := messages.LoginRequest{
		LoginToken: TestToken,
	}
	body, _ := json.Marshal(login)

	rr := apiTest(t, SimpleHandler(s.Login), "POST", body)
	assert.Equal(t, rr.Code, http.StatusOK)

	resp := messages.LoginResponse{}
	assert.Nil(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, resp.UserProfile.ID, TestUserID)

	d.AssertExpectations(t)
	m.AssertExpectations(t)
}
