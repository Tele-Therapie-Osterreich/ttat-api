package model

// TODO: REPLACE ALL THIS WITH AN IMAGE SERVICE.

import (
	"encoding/base64"
	"fmt"

	"github.com/pkg/errors"
)

// Image is a database model representing an avatar image associated
// with a therapist.
type Image struct {
	// Unique ID of the image.
	ID int `db:"id"`

	// Unique ID of the therapist the image is associated with.
	TherapistID int `db:"therapist_id"`

	// File extension for the image ("png" or "jpg").
	Extension string `db:"extension"`

	// Image data.
	Data []byte `db:"data"`
}

// MakeLink generates the URL required to access an image via the
// REST API.
func (image *Image) MakeLink() *string {
	s := fmt.Sprintf("/image/%d.%s", image.ID, image.Extension)
	return &s
}

// ImagePatch is a representation of an update to a therapist's image,
// giving the file extension ("jpg" or "png") and the binary image
// data.
type ImagePatch struct {
	Extension string
	Data      []byte
}

// DecodeImagePatch decodes the JSON representation of of an image
// patch.
func DecodeImagePatch(updates map[string]interface{}, k string) (*ImagePatch, error) {
	iPatchBody, ok := updates[k]
	if !ok {
		return nil, nil
	}
	patchBody, ok := iPatchBody.(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid image patch for '" + k + "'")
	}
	iExtension, ok := patchBody["extension"]
	if !ok {
		return nil, errors.New("missing extension in image patch for '" + k + "'")
	}
	extension, ok := iExtension.(string)
	if !ok || (extension != "jpg" && extension != "png") {
		return nil, errors.New("invalid extension in image patch for '" + k + "'")
	}
	iData, ok := patchBody["data"]
	if !ok {
		return nil, errors.New("missing data in image patch for '" + k + "'")
	}
	data, ok := iData.(string)
	if !ok {
		return nil, errors.New("invalid data in image patch for '" + k + "'")
	}
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	imagePatch := ImagePatch{Extension: extension, Data: decodedData}
	return &imagePatch, nil
}
