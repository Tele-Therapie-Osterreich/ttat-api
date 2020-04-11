package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Tele-Therapie-Osterreich/ttat-api/messages"
	"github.com/Tele-Therapie-Osterreich/ttat-api/mocks"
	"github.com/Tele-Therapie-Osterreich/ttat-api/model"
)

const TestEmail = "test@example.com"
const TestToken = "test-token"
const TestUserID = 1
const TestSession = "test-session"

func TestHealth(t *testing.T) {
	rr := apiTest(t, Health, "GET", []byte(""))

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, rr.Body.String(), `OK`)
}

func TestRequestLoginBadJSON(t *testing.T) {

}

func TestRequestLoginBadEmail(t *testing.T) {

}

func TestRequestLoginEmail(t *testing.T) {
	s, db, mailer := MockServer()

	db.On("CreateLoginToken", TestEmail, "en").Return(TestToken, nil)

	mailer.On("Send", "login-email-request", TestEmail, "en",
		map[string]string{"login_token": TestToken}).Return(nil)

	reqLoginEmail := messages.ReqLoginEmailRequest{
		Email:    TestEmail,
		Language: "en",
	}
	body, _ := json.Marshal(reqLoginEmail)

	rr := apiTest(t, SimpleHandler(s.RequestLoginEmail), "POST", body)

	assert.Equal(t, rr.Code, http.StatusNoContent)
	assert.Empty(t, rr.Body.String())
	mailer.AssertExpectations(t)
	db.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
	WithMocks(func(s *Server, db *mocks.DB, m *mocks.Mailer) {
		db.On("CheckLoginToken", TestToken).Return(TestEmail, "en", nil)
		u := model.User{
			ID:    TestUserID,
			Email: TestEmail,
		}
		db.On("LoginUser", TestEmail).Return(&u, nil, true, nil)
		db.On("CreateSession", TestUserID).Return(TestSession, nil)

		login := messages.LoginRequest{
			LoginToken: TestToken,
		}
		body, _ := json.Marshal(login)

		rr := apiTest(t, SimpleHandler(s.Login), "POST", body)
		assert.Equal(t, rr.Code, http.StatusOK)

		resp := messages.LoginResponse{}
		fmt.Println(rr.Body.String())
		assert.Nil(t, json.Unmarshal(rr.Body.Bytes(), &resp))
		assert.Equal(t, resp.UserProfile.ID, TestUserID)
	})
}
