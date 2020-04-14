package server

import (
	"github.com/Tele-Therapie-Osterreich/ttat-api/db"
	"github.com/go-chi/chi"
)

func (s *Server) routes(devMode bool, csrfSecret string,
	corsOrigins []string, testMode bool) chi.Router {
	r := chi.NewRouter()

	// Add common middleware.
	if !testMode {
		s.addMiddleware(r, devMode, csrfSecret, corsOrigins)
	}

	// Service health checks.
	r.Get("/", Health)

	// These routes need to be outside of the following block to exempt
	// them from CSRF protection since they're used before a session is
	// established.
	r.Post("/auth/request-login-email", SimpleHandler(s.requestLoginEmail))
	r.Post("/auth/login", SimpleHandler(s.login))

	// r.Get("/image/{id_and_extension}", SimpleHandler(s.imageDetail))

	r.Get("/therapist/{id:[0-9]+}", SimpleHandler(s.therapistDetail))

	r.Group(func(r chi.Router) {
		r.Use(CredentialCtx(s))

		// Authentication.
		r.Post("/auth/logout", SimpleHandler(s.logout))

		// Routes for authenticated user.
		r.Route("/me", s.selfRoutes)
	})

	return r
}

// Routes for viewing and manipulating user data.
func (s *Server) selfRoutes(r chi.Router) {
	// TODO: DECIDE EXACTLY HOW THIS SHOULD WORK
	r.Get("/", SimpleHandler(s.selfDetail(db.PreferPublic)))
	r.Get("/public", SimpleHandler(s.selfDetail(db.PublicOnly)))
	r.Get("/pending", SimpleHandler(s.selfDetail(db.PendingOnly)))
	// r.Patch("/", SimpleHandler(s.selfUpdate))
	r.Delete("/", SimpleHandler(s.selfDelete))
}
