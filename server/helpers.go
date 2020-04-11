package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tele-Therapie-Osterreich/ttat-api/mocks"
)

func apiTest(t *testing.T, handler http.HandlerFunc,
	method string, body []byte) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, "", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler(rr, req)
	return rr
}

func MockServer() (*Server, *mocks.DB, *mocks.Mailer) {
	m := mocks.Mailer{}
	db := mocks.DB{}
	return &Server{
		db:     &db,
		mailer: &m,
	}, &db, &m
}

func WithMocks(f func(s *Server, db *mocks.DB, m *mocks.Mailer)) {
	s, db, m := MockServer()
	f(s, db, m)
}
