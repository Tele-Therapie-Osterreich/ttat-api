package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"

	"github.com/Tele-Therapie-Osterreich/ttat-api/db"
	"github.com/Tele-Therapie-Osterreich/ttat-api/messages"
)

func (s *Server) imageDetail(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	imageIDAndExt := strings.Split(chi.URLParam(r, "id_and_extension"), ".")
	imageID, err := strconv.Atoi(imageIDAndExt[0])
	if len(imageIDAndExt) != 2 || err != nil {
		return BadRequest(w, messages.APIError{
			Code:    messages.InvalidImageFilename,
			Message: "invalid image filename for image detail",
		})
	}

	image, err := s.db.ImageByID(imageID)
	if err == db.ErrImageNotFound || imageIDAndExt[1] != image.Extension {
		return NotFound(w)
	}
	if err != nil {
		return nil, err
	}

	mimeType := ""
	switch image.Extension {
	case "jpg":
		mimeType = "image/jpeg"
	case "png":
		mimeType = "image/png"
	default:
		return NotFound(w)
	}
	w.Header().Add("Content-Type", mimeType)
	w.Write(image.Data)
	return nil, nil
}
