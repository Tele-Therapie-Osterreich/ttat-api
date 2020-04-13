package server

import (
	"net/http"

	"github.com/Tele-Therapie-Osterreich/ttat-api/db"
	"github.com/Tele-Therapie-Osterreich/ttat-api/messages"
	"github.com/Tele-Therapie-Osterreich/ttat-api/model"
	"github.com/Tele-Therapie-Osterreich/ttat-api/model/types"
)

// TODO: REFACTOR WITH selfDetail
func (s *Server) userDetail(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	userID := getIntParam(r, "id")
	if userID == nil {
		return NotFound(w)
	}

	user, err := s.db.UserByID(*userID)
	if err == db.ErrUserNotFound {
		return NotFound(w)
	}
	if err != nil {
		return nil, err
	}
	if user.Status != types.Approved {
		return NotFound(w)
	}

	image, err := s.db.ImageByUserID(*userID)
	if err != db.ErrUserNotFound && err != nil {
		return nil, err
	}

	// TODO: Also need specialities...
	return model.UserFullProfileFromUser(user, image), nil
}

func (s *Server) selfDetail(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	userID := accessControl(r, true)
	if userID == nil {
		return NotFound(w)
	}

	user, err := s.db.UserByID(*userID)
	if err == db.ErrUserNotFound {
		return NotFound(w)
	}
	if err != nil {
		return nil, err
	}

	image, err := s.db.ImageByUserID(*userID)
	if err != db.ErrUserNotFound && err != nil {
		return nil, err
	}

	// TODO: Also need specialities...
	return model.UserFullProfileFromUser(user, image), nil
}

func (s *Server) selfDelete(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	userID := accessControl(r, false)
	if userID == nil {
		return NotFound(w)
	}

	err := s.db.DeleteUser(*userID)
	if err == db.ErrUserNotFound {
		return NotFound(w)
	}
	if err != nil {
		return nil, err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil, nil
}

func (s *Server) selfUpdate(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	userID := accessControl(r, false)
	if userID == nil {
		return NotFound(w)
	}

	// Read patch request body.
	body, err := ReadBody(r, 0)
	if err != nil {
		return BadRequest(w, messages.APIError{
			Code:    messages.InvalidJSONBody,
			Message: "invalid JSON body for therapist account update",
		})
	}

	// Look up user value and patch it.
	user, err := s.db.UserByID(*userID)
	if err != nil {
		return nil, err
	}
	image, err := s.db.ImageByUserID(*userID)
	if err != nil {
		return nil, err
	}
	imagePatch, err := user.Patch(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, nil
	}
	if err = s.db.UpdateUser(user); err != nil {
		return nil, err
	}

	// Deal with image patches.
	if imagePatch != nil {
		if image == nil {
			newImage := model.Image{UserID: *userID}
			image = &newImage
		}
		image.Extension = imagePatch.Extension
		image.Data = imagePatch.Data
		if image, err = s.db.UpsertImage(image); err != nil {
			return nil, err
		}
	}

	return model.UserFullProfileFromUser(user, image), nil
}
