package server

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/Tele-Therapie-Osterreich/ttat-api/db"
	"github.com/Tele-Therapie-Osterreich/ttat-api/mail"
)

// Server is the server structure for the TTAT API service.
type Server struct {
	srv           *http.Server
	db            db.DB
	secureSession bool
	mailer        mail.Mailer
}

// NewServer creates the server structure for the TTAT API service.
func NewServer(cfg *Config, db db.DB, mailer mail.Mailer) *Server {
	// Fixed CORS origin list from environment.
	corsOrigins := []string{}
	if len(cfg.CORSOrigins) > 0 {
		corsOrigins = strings.Split(cfg.CORSOrigins, ",")
	}

	// Randomise ID generation.
	rand.Seed(int64(time.Now().Nanosecond()))

	// Basic server initialisation.
	s := &Server{
		secureSession: !cfg.DevMode,
		db:            db,
		mailer:        mailer,
	}
	s.srv = &http.Server{
		Handler: s.routes(cfg.DevMode, cfg.CSRFSecret, corsOrigins),
		Addr:    fmt.Sprintf(":%d", cfg.Port),
	}

	return s
}

// Serve runs a server event loop.
func (s *Server) Serve() {
	errChan := make(chan error, 0)
	go func() {
		log.Info().
			Str("address", s.srv.Addr).
			Msg("server started")
		err := s.srv.ListenAndServe()
		if err != nil {
			errChan <- err
		}
	}()

	signalCh := make(chan os.Signal, 0)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	var err error

	select {
	case <-signalCh:
	case err = <-errChan:
	}

	s.shutdown()

	if err == nil {
		log.Info().Msg("server shutting down")
	} else {
		log.Fatal().Err(err).Msg("server failed")
	}
}

// Shut down server.
func (s *Server) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Error().Err(err)
	}
}
