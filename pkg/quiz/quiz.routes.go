package quiz

import (
	"Back-end/pkg/security"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterQuizRoutes(router chi.Router) {
	router.Route("/quiz", func(r chi.Router) {
		r.Use(middleware.Logger)
		// Routes non protégées
		r.Get("/", GetAllQuizzes)
		r.Get("/byQuery", GetQuizByID)
		// Routes protégées
		r.With(security.Middleware).Post("/create", CreateQuiz)
		r.With(security.Middleware).Put("/byQuery", UpdateQuiz)
		r.With(security.Middleware).Delete("/byQuery", DeleteQuiz)
	})
}
