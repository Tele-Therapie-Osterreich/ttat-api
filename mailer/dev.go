package mailer

import (
	"fmt"
)

// DevMailer is a development email sender.
type DevMailer struct{}

// NewDevMailer creates a new development email sender.
func NewDevMailer() *DevMailer {
	return &DevMailer{}
}

// Send sends an email.
func (m *DevMailer) Send(template string, language string,
	email string, data map[string]string) error {
	fmt.Println("====> EMAIL SEND ->", email)
	fmt.Println("  template =", template, "   language =", language)
	for k, v := range data {
		fmt.Println(" ", k, "=", v)
	}
	fmt.Println("<==== EMAIL SEND")
	return nil
}
