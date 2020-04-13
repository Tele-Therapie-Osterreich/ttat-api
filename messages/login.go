package messages

import "github.com/Tele-Therapie-Osterreich/ttat-api/model"

// ReqLoginEmailRequest is a message structure for requests to send an
// email containing a login token.
type ReqLoginEmailRequest struct {
	Email    string `json:"email"`
	Language string `json:"language"`
}

// LoginRequest is a message structure for login requests.
type LoginRequest struct {
	LoginToken string `json:"login_token"`
}

// LoginResponse is a message structure for successful login
// responses.
type LoginResponse struct {
	Profile      *model.TherapistFullProfile `json:"profile"`
	NewTherapist bool                        `json:"new_therapist"`
}
