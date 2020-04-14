package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tele-Therapie-Osterreich/ttat-api/mocks"
	"github.com/go-chi/chi"
)

// Test setup constants.
const (
	TestEmail          = "test@example.com"
	TestName           = "Test User"
	TestToken          = "test-token"
	TestTherapistID    = 2345
	TestProfileID      = 6789
	TestImageID        = 997
	TestImageExtension = "jpg"
	TestSession        = "test-session"
)

func testTherapistURL(t string) string {
	return fmt.Sprintf(t, TestTherapistID)
}

func testImageURL() string {
	return fmt.Sprintf("/image/%d.%s", TestImageID, TestImageExtension)
}

type apiOptions struct {
	body     interface{}
	bodyJSON []byte
	session  string
}

func apiTest(t *testing.T, r chi.Router,
	method string, url string, opts *apiOptions) *httptest.ResponseRecorder {
	body := []byte{}
	if opts != nil {
		body = opts.bodyJSON
		if opts.body != nil {
			var err error
			body, err = json.Marshal(opts.body)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if opts != nil && opts.session != "" {
		req.AddCookie(&http.Cookie{Name: "session", Value: opts.session})
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func mockServer() (*mocks.DB, *mocks.Mailer, chi.Router) {
	m := mocks.Mailer{}
	db := mocks.DB{}
	s := Server{
		db:     &db,
		mailer: &m,
	}
	return &db, &m, s.routes(true, "test-secret", []string{}, true)
}
