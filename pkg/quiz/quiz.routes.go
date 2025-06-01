package quiz

import (
	"Back-end/pkg/quiz/controller"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterQuizRoutes(router chi.Router) {
	router.Route("/quiz", func(r chi.Router) {
		r.Use(middleware.Logger)

		// ✅ Routes fonctionnelles
		r.Get("/", controller.GetAllQuizzes)
		r.Post("/create", controller.CreateQuiz)

		// ✅ Nouvelles routes avec ID en query param
		router.Get("/quizByQuery", controller.GetQuizByID)
		router.Put("/quizByQuery", controller.UpdateQuiz)
		router.Delete("/quizByQuery", controller.DeleteQuiz)
	})
	log.Println("✅ Routes quiz enregistrées !")
}
