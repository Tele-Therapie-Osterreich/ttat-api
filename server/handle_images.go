package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Tele-Therapie-Osterreich/ttat-api/chassis"
	"github.com/Tele-Therapie-Osterreich/ttat-api/db"
	"github.com/go-chi/chi"
)

func (s *Server) imageDetail(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	imageIDAndExt := strings.Split(chi.URLParam(r, "id_and_extension"), ".")
	imageID, err := strconv.Atoi(imageIDAndExt[0])
	if len(imageIDAndExt) != 2 || err != nil {
		return chassis.BadRequest(w, "invalid image filename")
	}

	image, err := s.db.ImageByID(imageID)
	if err == db.ErrImageNotFound || imageIDAndExt[1] != image.Extension {
		return chassis.NotFound(w)
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
		return chassis.NotFound(w)
	}
	w.Header().Add("Content-Type", mimeType)
	w.Write(image.Data)
	return nil, nil
}