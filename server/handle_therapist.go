package server

import (
	"net/http"

	"github.com/Tele-Therapie-Osterreich/ttat-api/db"
	"github.com/Tele-Therapie-Osterreich/ttat-api/model"
	"github.com/Tele-Therapie-Osterreich/ttat-api/model/types"
)

// TODO: REFACTOR WITH selfDetail
func (s *Server) therapistDetail(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	thID := getIntParam(r, "id")
	if thID == nil {
		return NotFound(w)
	}

	info, err := s.db.TherapistInfoByID(*thID, db.PublicOnly)
	if err == db.ErrTherapistNotFound {
		return NotFound(w)
	}
	if err != nil {
		return nil, err
	}
	if info.Base.Status != types.Active {
		return NotFound(w)
	}

	// TODO: ALSO NEED SPECIALITIES AND IMAGE HANDLING...
	return model.NewTherapistFullView(info), nil
}

func (s *Server) selfDetail(profile db.ProfileSelection) SimpleHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		thID := accessControl(r)
		if thID == nil {
			return NotFound(w)
		}

		info, err := s.db.TherapistInfoByID(*thID, profile)
		if err == db.ErrTherapistNotFound {
			return NotFound(w)
		}
		if err != nil {
			return nil, err
		}

		// TODO: Also need specialities...
		return model.NewTherapistFullView(info), nil
	}
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

// func (s *Server) selfUpdate(w http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	thID := accessControl(r)
// 	if thID == nil {
// 		return NotFound(w)
// 	}

// 	// Read patch request body.
// 	body, err := ReadBody(r, 0)
// 	if err != nil {
// 		return BadRequest(w, messages.APIError{
// 			Code:    messages.InvalidJSONBody,
// 			Message: "invalid JSON body for therapist account update",
// 		})
// 	}

// 	// Look up therapist profile and patch it.
// 	imagePatch, err := s.db.UpdateTherapistProfile(*thID, body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Deal with image patches.
// 	if imagePatch != nil {
// 		// if image == nil {
// 		// 	newImage := model.Image{TherapistID: *thID}
// 		// 	image = &newImage
// 		// }
// 		// image.Extension = imagePatch.Extension
// 		// image.Data = imagePatch.Data
// 		// if image, err = s.db.UpsertImage(image); err != nil {
// 		// 	return nil, err
// 		// }
// 	}

// 	return model.NewTherapistFullView(th, image), nil
// }
