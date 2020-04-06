package mailer

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

// MailEvent is a structure containing information about an email to
// be sent.
type MailEvent struct {
	Template string
	Email    string
	Language string
	Data     map[string]string
}

// MailSender is the main email sender goroutine: runs off of
// multiplexed message channel.
func MailSender(mailCh chan MailEvent, mailer Mailer) {
	// TODO: HANDLE MESSAGES IN PARALLEL WITH RATE LIMITING.
	for {
		ev := <-mailCh

		language := ev.Language
		if language == "" {
			language = "en"
		}
		err := mailer.Send(ev.Template, language, ev.Email, ev.Data)
		if err != nil {
			log.Error().Err(err).
				Str("template", ev.Template).
				Str("email", ev.Email).
				Str("data", fmt.Sprintf("%v", ev.Data)).
				Msg("mailer couldn't send email")
			continue
		}
		log.Info().
			Str("template", ev.Template).
			Str("email", ev.Email).
			Str("data", fmt.Sprintf("%v", ev.Data)).
			Msg("email sent")
	}
}
