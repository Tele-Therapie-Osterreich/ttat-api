package server

import (
	"github.com/go-chi/chi"

	"github.com/Tele-Therapie-Osterreich/ttat-api/chassis"
)

func (s *Server) routes(devMode bool, csrfSecret string, corsOrigins []string) chi.Router {
	r := chi.NewRouter()

	// Add common middleware.
	s.addMiddleware(r, devMode, csrfSecret, corsOrigins)

	// Service health checks.
	r.Get("/", chassis.Health)

	// These routes need to be outside of the following block to exempt
	// them from CSRF protection since they're used before a session is
	// established.
	r.Post("/auth/request-login-email", chassis.SimpleHandler(s.requestLoginEmail))
	r.Post("/auth/login", chassis.SimpleHandler(s.login))

	r.Get("/image/{id_and_extension}", chassis.SimpleHandler(s.imageDetail))

	r.Group(func(r chi.Router) {
		r.Use(CredentialCtx(s))

		// Authentication.
		r.Post("/auth/logout", chassis.SimpleHandler(s.logout))

		// Routes for authenticated user.
		r.Route("/me", s.userRoutes)
	})

	return r
}

// Routes for viewing and manipulating user data.
func (s *Server) userRoutes(r chi.Router) {
	r.Get("/", chassis.SimpleHandler(s.userDetail))
	r.Patch("/", chassis.SimpleHandler(s.userUpdate))
	r.Delete("/", chassis.SimpleHandler(s.userDelete))
}
