package proposition

import (
	"Back-end/pkg/security"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterPropositionRoutes(router chi.Router) {
	router.Route("/proposition", func(r chi.Router) {
		r.Use(middleware.Logger)
		// Routes non protégées
		r.Get("/all", GetAllPropositions)
		r.Get("/get", GetPropositionByID)
		// Routes protégées
		r.With(security.Middleware).Post("/create", CreateProposition)
		r.With(security.Middleware).Put("/update", UpdateProposition)
		r.With(security.Middleware).Delete("/delete", DeleteProposition)
	})
}
