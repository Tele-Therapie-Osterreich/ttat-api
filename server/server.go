package server

import (
	"context"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/Tele-Therapie-Osterreich/ttat-api/chassis"
	"github.com/Tele-Therapie-Osterreich/ttat-api/db"
	"github.com/Tele-Therapie-Osterreich/ttat-api/mailer"
)

// Server is the server structure for the user service.
type Server struct {
	chassis.Server
	db            db.DB
	secureSession bool
	mailer        mailer.Mailer
	mailCh        chan mailer.MailEvent
}

// Config contains the configuration information needed to start
// the user service.
type Config struct {
	DevMode            bool   `env:"DEV_MODE,default=false"`
	DBURL              string `env:"DATABASE_URL,required"`
	Port               int    `env:"PORT,default=8080"`
	CSRFSecret         string `env:"CSRF_SECRET"`
	CORSOrigins        string `env:"CORS_ORIGINS"`
	MJPublicKey        string `env:"MAILJET_API_KEY_PUBLIC"`
	MJPrivateKey       string `env:"MAILJET_API_KEY_PRIVATE"`
	SimultaneousEmails int    `env:"SIMULTANEOUS_EMAILS,default=10"`
}

// NewServer creates the server structure for the user service.
func NewServer(cfg *Config) *Server {
	// Fixed CORS origin list from environment.
	corsOrigins := []string{}
	if len(cfg.CORSOrigins) > 0 {
		corsOrigins = strings.Split(cfg.CORSOrigins, ",")
	}

	// Common server initialisation.
	s := &Server{
		secureSession: !cfg.DevMode,
	}
	s.Init(cfg.Port, s.routes(cfg.DevMode, cfg.CSRFSecret, corsOrigins))

	// Connect to database.
	timeout, _ := context.WithTimeout(context.Background(), time.Second*10)
	var err error
	s.db, err = db.NewPGClient(timeout, cfg.DBURL)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't connect to user database")
	}

	// Initialise mailer.
	if cfg.MJPublicKey == "" || cfg.MJPrivateKey == "" ||
		cfg.MJPublicKey == "dev" || cfg.MJPrivateKey == "dev" {
		log.Info().Msg("using development mailer")
		s.mailer = mailer.NewDevMailer()
	} else {
		s.mailer, err = mailer.NewMailjetMailer(cfg.MJPublicKey, cfg.MJPrivateKey)
		if err != nil {
			log.Fatal().Err(err).Msg("couldn't connect to Mailjet")
		}
	}

	// Set up mail sender.
	s.mailCh = make(chan mailer.MailEvent, cfg.SimultaneousEmails)
	go mailer.MailSender(s.mailCh, s.mailer)

	return s
}

// SendEmail is a simple interface for sending emails.
func (s *Server) SendEmail(template string, email string, language string,
	data map[string]string) {
	s.mailCh <- mailer.MailEvent{
		Template: template,
		Email:    email,
		Language: language,
		Data:     data,
	}
}
