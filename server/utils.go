package server

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

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
