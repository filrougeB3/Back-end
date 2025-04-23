package quiz

import (
	"Back-end/pkg/quiz/controller"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterQuizRoutes(router chi.Router) {
	router.Route("/quiz", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Get("/", controller.GetAllQuizzes)
		r.Get("/{id}", controller.GetQuizByID)
		r.Post("/create", controller.CreateQuiz)
		r.Put("/{id}", controller.UpdateQuiz)
		r.Delete("/{id}", controller.DeleteQuiz)
	})
}
