package server

import (
	"encoding/json"
	"net/http"

	"github.com/Tele-Therapie-Osterreich/ttat-api/messages"
)

// BadRequest sets up an HTTP 400 Bad Request with a given error
// message and returns the (nil, nil) pair used by SimpleHandler to
// signal that the response has been dealt with.
func BadRequest(w http.ResponseWriter, err messages.APIError) (interface{}, error) {
	body, _ := json.Marshal(err)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(body)
	return nil, nil
}

// NotFound sets up an HTTP 404 Not Found and returns the (nil, nil)
// pair used by SimpleHandler to signal that the response has been
// dealt with.
func NotFound(w http.ResponseWriter) (interface{}, error) {
	http.NotFound(w, nil)
	return nil, nil
}

// Forbidden sets up an HTTP 403 Forbidden and returns the (nil, nil)
// pair used by SimpleHandler to signal that the response has been
// dealt with.
func Forbidden(w http.ResponseWriter) (interface{}, error) {
	w.WriteHeader(http.StatusForbidden)
	return nil, nil
}

// NoContent sets up an HTTP 204 No Content and returns the (nil, nil)
// pair used by SimpleHandler to signal that the response has been
// dealt with.
func NoContent(w http.ResponseWriter) (interface{}, error) {
	w.WriteHeader(http.StatusNoContent)
	return nil, nil
}
