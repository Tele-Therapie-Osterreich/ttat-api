package mail

import (
	"errors"
)

// ErrUnknownEmailTemplate is the error returned by a template store
// when an unknown template is requested.
var ErrUnknownEmailTemplate = errors.New("email template unknown")

// Mailer represents machinery for sending template-based emails.
type Mailer interface {
	Send(template string, language string, email string, data map[string]string)
}
