package server

import (
	"net/http"

	"github.com/Tele-Therapie-Osterreich/ttat-api/chassis"
	"github.com/Tele-Therapie-Osterreich/ttat-api/db"
	"github.com/Tele-Therapie-Osterreich/ttat-api/model"
)

func (s *Server) userDetail(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	userID, _ := accessControl(r)
	if userID == nil {
		return chassis.NotFound(w)
	}

	user, err := s.db.UserByID(*userID)
	if err == db.ErrUserNotFound {
		return chassis.NotFound(w)
	}
	if err != nil {
		return nil, err
	}

	image, err := s.db.ImageByUserID(*userID)
	if err != db.ErrUserNotFound && err != nil {
		return nil, err
	}

	return model.UserFullProfileFromUser(user, image), nil
}

func (s *Server) userDelete(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	userID, _ := accessControl(r)
	if userID == nil {
		return chassis.NotFound(w)
	}

	err := s.db.DeleteUser(*userID)
	if err == db.ErrUserNotFound {
		return chassis.NotFound(w)
	}
	if err != nil {
		return nil, err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil, nil
}

func (s *Server) userUpdate(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// userID, actingUserID, adminAction := accessControl(r)
	userID, actingUserID := accessControl(r)
	if userID == nil {
		return chassis.NotFound(w)
	}
	if *actingUserID != *userID { // && !adminAction {
		http.Error(w, "illegal patch", http.StatusForbidden)
		return nil, nil
	}

	// Read patch request body.
	body, err := chassis.ReadBody(r, 0)
	if err != nil {
		return chassis.BadRequest(w, err.Error())
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
