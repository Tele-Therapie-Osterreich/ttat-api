package messages

// APIErrorCode is a short code used for API error messages.
type APIErrorCode string

// API error message codes.
const (
	InvalidJSONBody      APIErrorCode = "invalid-json-body"
	InvalidEmailAddress  APIErrorCode = "invalid-email-address"
	UnknownLoginToken    APIErrorCode = "unknown-login-token"
	InvalidImageFilename APIErrorCode = "invalid-image-filename"
)

// APIError is a response message structure for API error messages.
type APIError struct {
	Code    APIErrorCode `json:"code"`
	Message string       `json:"message"`
	Field   string       `json:"field",omitempty`
}
