package auth

import (
	"github.com/go-chi/chi/v5"
)

func RegisterAuthRoutes(router chi.Router) {
	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", CreateUser)
		r.Post("/login", LoginUser)
	})
}
