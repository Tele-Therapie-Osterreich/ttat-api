package mail

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"time"

	mailjet "github.com/mailjet/mailjet-apiv3-go"
	"github.com/mailjet/mailjet-apiv3-go/resources"
	"github.com/rs/zerolog/log"
)

var fromRE = regexp.MustCompile(`(.+) <([^>]+)>`)

// MailjetMailer sends email using Mailjet.
type MailjetMailer struct {
	mailCh chan event
	mj     *mailjet.Client
	mu     sync.Mutex
	tmpls  map[string]templateInfo
}

// Fields we need to extract from the template definition to use the
// sending API.
type templateInfo struct {
	id          int64
	fromEmail   string
	fromName    string
	replyTo     string
	replyEmail  string
	senderEmail string
	senderName  string
	subject     string
}

// NewMailjetMailer creates a new mailer for Mailjet based on Mailjet
// API key credentials.
func NewMailjetMailer(pubkey, privkey string, simultaneousEmails int) (*MailjetMailer, error) {
	// Set up mail sender.
	mailCh := make(chan event, simultaneousEmails)

	mailer := MailjetMailer{}
	log.Info().Msg("connecting to Mailjet")
	mailer.mj = mailjet.NewMailjetClient(pubkey, privkey)
	err := mailer.loadTemplates()
	if err != nil {
		return nil, err
	}
	go mailer.templateUpdater()
	go mailer.sender(mailCh)
	return &mailer, nil
}

// A structure containing information about an email to be sent.
type event struct {
	Template string
	Email    string
	Language string
	Data     map[string]string
}

// Main email sender goroutine: runs off of multiplexed message
// channel.
func (m *MailjetMailer) sender(mailCh chan event) {
	// TODO: HANDLE MESSAGES IN PARALLEL WITH RATE LIMITING.
	for {
		ev := <-mailCh

		language := ev.Language
		if language == "" {
			language = "en"
		}
		err := m.doSend(ev.Template, language, ev.Email, ev.Data)
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

// loadTemplates loads all our template from Mailjet.
func (m *MailjetMailer) loadTemplates() error {
	// Get basic template information first.
	var tmpls []resources.Template
	_, _, err := m.mj.List("template", &tmpls, mailjet.Filter("OwnerType", "user"))
	if err != nil {
		log.Error().Err(err).Msg("failed to load templates from Mailjet")
		return err
	}

	// Now look up the detailed information for each template (have to
	// do these one at a time).
	tmplInfo := map[string]templateInfo{}
	for _, tmpl := range tmpls {
		req := mailjet.Request{
			Resource: "template",
			ID:       tmpl.ID,
			Action:   "detailcontent",
		}

		// Get detailed template information and make sure it matches what
		// we expect to see.
		var res []resources.TemplateDetailcontent
		err = m.mj.Get(&req, &res)
		if err != nil {
			log.Error().Err(err).
				Str("name", tmpl.Name).
				Msg("getting template details")
			continue
		}
		if len(res) != 1 {
			log.Error().
				Str("name", tmpl.Name).
				Msg("strange reply from Mailjet for template details")
			continue
		}
		headerMap, ok := res[0].Headers.(map[string]interface{})
		if !ok {
			log.Error().
				Str("name", tmpl.Name).
				Msg("invalid headers in template details from Mailjet")
			continue
		}

		// The "From" header inconveniently includes the From email and
		// name together. We need them separately for sending mail.
		from := fromRE.FindStringSubmatch(headerMap["From"].(string))

		// Collect template information to save for later use when sending
		// mail.
		tmplInfo[tmpl.Name] = templateInfo{
			id:          tmpl.ID,
			fromEmail:   from[2],
			fromName:    from[1],
			replyTo:     headerMap["Reply-To"].(string),
			replyEmail:  headerMap["ReplyEmail"].(string),
			senderEmail: headerMap["SenderEmail"].(string),
			senderName:  headerMap["SenderName"].(string),
			subject:     headerMap["Subject"].(string),
		}
	}

	// Save the template information.
	m.mu.Lock()
	defer m.mu.Unlock()
	oldtmpls := m.tmpls
	m.tmpls = tmplInfo
	for name, tmpl := range m.tmpls {
		old, chk := oldtmpls[name]
		if !chk {
			log.Info().
				Str("name", name).
				Int64("id", tmpl.id).
				Msg("loaded new Mailjet template")
			continue
		}
		if old != tmpl {
			log.Info().
				Str("name", name).
				Int64("id", tmpl.id).
				Msg("updated existing Mailjet template")
		}
	}
	return nil
}

// Repeatedly reload templates.
func (m *MailjetMailer) templateUpdater() {
	for range time.Tick(3 * time.Minute) {
		m.loadTemplates()
	}
}

// Send queues an email for sending.
func (m *MailjetMailer) Send(template string, email string, language string,
	data map[string]string) {
	m.mailCh <- event{
		Template: template,
		Email:    email,
		Language: language,
		Data:     data,
	}
}

// Send sends an email using a Mailjet template.
func (m *MailjetMailer) doSend(template string, language string,
	email string, data map[string]string) error {
	// TODO: MULTI-LINGUAL TEMPLATES

	// Get Mailjet template.
	tmplInfo, ok := m.getTemplate(template)
	if !ok {
		return ErrUnknownEmailTemplate
	}

	// Build mail information and send mail.
	info := &mailjet.InfoSendMail{
		FromEmail: tmplInfo.fromEmail,
		FromName:  tmplInfo.fromName,
		Recipients: []mailjet.Recipient{
			{
				Email: email,
			},
		},
		Subject:            tmplInfo.subject,
		MjTemplateLanguage: "true",
		MjTemplateID:       strconv.Itoa(int(tmplInfo.id)),
		Vars:               data,
	}
	res, err := m.mj.SendMail(info)
	if err != nil {
		return err
	}
	if len(res.Sent) != 1 {
		err = errors.New("invalid result from Mailjet")
		log.Error().Err(err)
		return err
	}

	// Log mail send.
	log.Info().
		Int64("mailjet-message-id", res.Sent[0].MessageID).
		Str("template", template).
		Msg("email sent")
	return nil
}

// Look up a template ID from its name in our map.
func (m *MailjetMailer) getTemplate(name string) (*templateInfo, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.tmpls[name]
	return &t, ok
}
