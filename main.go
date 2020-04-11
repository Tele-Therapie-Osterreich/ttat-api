package main

import (
	"context"
	"os"
	"time"

	"github.com/joeshaw/envdecode"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Tele-Therapie-Osterreich/ttat-api/db"
	"github.com/Tele-Therapie-Osterreich/ttat-api/mail"
	"github.com/Tele-Therapie-Osterreich/ttat-api/server"
)

func main() {
	// Read configuration from environment variables.
	cfg := server.Config{}
	err := envdecode.StrictDecode(&cfg)
	if err != nil {
		log.Fatal().Err(err).
			Msg("failed to process environment variables")
	}

	// Set up logging.
	baselog := zerolog.New(os.Stdout)
	if cfg.DevMode {
		baselog = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})
	}
	applog := baselog.With().Timestamp().Logger()
	log.Logger = applog

	// Connect to database.
	timeout, _ := context.WithTimeout(context.Background(), time.Second*10)
	db, err := db.NewPGClient(timeout, cfg.DBURL)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't connect to database")
	}

	// Initialise mailer.
	var mailer mail.Mailer
	if cfg.MJPublicKey == "" || cfg.MJPrivateKey == "" ||
		cfg.MJPublicKey == "dev" || cfg.MJPrivateKey == "dev" {
		log.Info().Msg("using development mailer")
		mailer = mail.NewDevMailer()
	} else {
		var err error
		mailer, err = mail.NewMailjetMailer(cfg.MJPublicKey, cfg.MJPrivateKey,
			cfg.SimultaneousEmails)
		if err != nil {
			log.Fatal().Err(err).Msg("couldn't connect to Mailjet")
		}
	}

	// Build and run server.
	serv := server.NewServer(&cfg, db, mailer)
	serv.Serve()
}
