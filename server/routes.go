package server

import (
	"github.com/go-chi/chi"
)

func (s *Server) routes(devMode bool, csrfSecret string, corsOrigins []string) chi.Router {
	r := chi.NewRouter()

	// Add common middleware.
	s.addMiddleware(r, devMode, csrfSecret, corsOrigins)

	// Service health checks.
	r.Get("/", Health)

	// These routes need to be outside of the following block to exempt
	// them from CSRF protection since they're used before a session is
	// established.
	r.Post("/auth/request-login-email", SimpleHandler(s.RequestLoginEmail))
	r.Post("/auth/login", SimpleHandler(s.Login))

	// r.Get("/image/{id_and_extension}", SimpleHandler(s.imageDetail))

	// r.Get("/user/{id:[0-9]+}", SimpleHandler(s.userDetail))

	// r.Group(func(r chi.Router) {
	// 	r.Use(CredentialCtx(s))

	// 	// Authentication.
	// 	r.Post("/auth/logout", SimpleHandler(s.logout))

	// 	// Routes for authenticated user.
	// 	r.Route("/me", s.userRoutes)
	// })

	return r
}

// // Routes for viewing and manipulating user data.
// func (s *Server) userRoutes(r chi.Router) {
// 	r.Get("/", SimpleHandler(s.userDetail))
// 	r.Patch("/", SimpleHandler(s.userUpdate))
// 	r.Delete("/", SimpleHandler(s.userDelete))
// }
