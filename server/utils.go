package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

// SimpleHandlerFunc is a HTTP handler function that signals internal
// errors by returning a normal Go error, and when successful returns
// a response body to be marshalled to JSON. It can be wrapped in the
// SimpleHandler middleware to produce a normal HTTP handler function.
type SimpleHandlerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error)

// SimpleHandler wraps a simpleHandler-style HTTP handler function to
// turn it into a normal HTTP handler function. Go errors from the
// inner handler are returned to the caller as "500 Internal Server
// Error" responses. Returns from successful processing by the inner
// handler and marshalled into a JSON response body.
func SimpleHandler(inner SimpleHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Run internal handler: returns a marshalable result and an
		// error, either of which may be nil.
		result, err := inner(w, r)

		// Propagate Go errors as "500 Internal Server Error" responses.
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("handling %q: %v", r.RequestURI, err)
			return
		}

		// No response body, so internal handler dealt with response
		// setup.
		if result == nil {
			return
		}

		// Marshal JSON response body.
		body, err := json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("handling %q: %v", r.RequestURI, err)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(body)
	}
}

func getIntParam(r *http.Request, p string) *int {
	params := chi.URLParam(r, p)
	if params == "" {
		return nil
	}
	param, err := strconv.Atoi(params)
	if err != nil {
		return nil
	}
	return &param
}
