package user

import "github.com/go-chi/chi/v5"

func RegisterUserRoutes(router chi.Router) {
	router.Route("/user", func(r chi.Router) {
		// Routes protégées
		r.With(UserMiddleware).Get("/", GetUser)
		r.With(UserMiddleware).Put("/", UpdateUser)
	})
}