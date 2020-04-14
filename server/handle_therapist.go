package server

import (
	"net/http"

	"github.com/Tele-Therapie-Osterreich/ttat-api/db"
	"github.com/Tele-Therapie-Osterreich/ttat-api/messages"
	"github.com/Tele-Therapie-Osterreich/ttat-api/model"
	"github.com/Tele-Therapie-Osterreich/ttat-api/model/types"
)

// TODO: REFACTOR WITH selfDetail
func (s *Server) therapistDetail(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	thID := getIntParam(r, "id")
	if thID == nil {
		return NotFound(w)
	}

	th, err := s.db.TherapistByID(*thID)
	if err == db.ErrTherapistNotFound {
		return NotFound(w)
	}
	if err != nil {
		return nil, err
	}
	if th.Status != types.Approved {
		return NotFound(w)
	}

	image, err := s.db.ImageByTherapistID(*thID)
	if err != db.ErrTherapistNotFound && err != nil {
		return nil, err
	}

	// TODO: Also need specialities...
	return model.TherapistFullProfileFromTherapist(th, image), nil
}

func (s *Server) selfDetail(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	thID := accessControl(r)
	if thID == nil {
		return NotFound(w)
	}

	th, err := s.db.TherapistByID(*thID)
	if err == db.ErrTherapistNotFound {
		return NotFound(w)
	}
	if err != nil {
		return nil, err
	}

	image, err := s.db.ImageByTherapistID(*thID)
	if err != db.ErrTherapistNotFound && err != nil {
		return nil, err
	}

	// TODO: Also need specialities...
	return model.TherapistFullProfileFromTherapist(th, image), nil
}

func (s *Server) selfDelete(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	thID := accessControl(r)
	if thID == nil {
		return NotFound(w)
	}

	err := s.db.DeleteTherapist(*thID)
	if err == db.ErrTherapistNotFound {
		return NotFound(w)
	}
	if err != nil {
		return nil, err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil, nil
}

func (s *Server) selfUpdate(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	thID := accessControl(r)
	if thID == nil {
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

	// Look up therapist value and patch it.

	edits, err := s.db.AddPendingEdits(*thID, body)
	if err != nil {
		return nil, err
	}

	// TODO: SHOULD BE GETTING THE THERAPIST DATA WITH ANY EXISTING
	// PATCHES APPLIED HERE.
	th, err := s.db.TherapistByID(*thID)
	if err != nil {
		return nil, err
	}
	image, err := s.db.ImageByTherapistID(*thID)
	if err != nil {
		return nil, err
	}
	imagePatch, err := th.Patch(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, nil
	}
	// if err = s.db.UpdateTherapist(th); err != nil {
	// 	return nil, err
	// }

	// Deal with image patches.
	if imagePatch != nil {
		if image == nil {
			newImage := model.Image{TherapistID: *thID}
			image = &newImage
		}
		image.Extension = imagePatch.Extension
		image.Data = imagePatch.Data
		if image, err = s.db.UpsertImage(image); err != nil {
			return nil, err
		}
	}

	return model.TherapistFullProfileFromTherapist(th, image), nil
}
