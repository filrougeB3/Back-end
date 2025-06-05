package question

import (
	"Back-end/pkg/security"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterQuestionRoutes(router chi.Router) {
	router.Route("/question", func(r chi.Router) {
		r.Use(middleware.Logger)
		// Routes non protégées
		r.Get("/all", GetAllQuestions)
		r.Get("/get", GetQuestionByID)
		// Routes protégées
		r.With(security.Middleware).Post("/create", CreateQuestion)
		r.With(security.Middleware).Put("/update", UpdateQuestion)
		r.With(security.Middleware).Delete("/delete", DeleteQuestion)
	})
}
