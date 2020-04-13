package server

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	_, _, r := mockServer()

	rr := apiTest(t, r, "GET", "/", nil)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, rr.Body.String(), `OK`)
}
