package user

import (
	"Back-end/pkg/security"

	"github.com/go-chi/chi/v5"
)

func RegisterUserRoutes(router chi.Router) {
	router.Route("/user", func(r chi.Router) {
		// Routes protégées
		r.With(security.Middleware).Get("/", GetUser)
		r.With(security.Middleware).Put("/", UpdateUser)
	})
}
