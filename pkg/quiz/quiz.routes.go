package quiz

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterQuizRoutes(router chi.Router) {
	router.Route("/quiz", func(r chi.Router) {
		r.Use(middleware.Logger)

		// ✅ Routes classiques
		r.Get("/", GetAllQuizzes)
		r.Post("/create", CreateQuiz)

		// ✅ Routes avec ID en query param : /quiz/byQuery?id=42
		r.Get("/byQuery", GetQuizByID)
		r.Put("/byQuery", UpdateQuiz)
		r.Delete("/byQuery", DeleteQuiz)
	})

	log.Println("✅ Routes quiz enregistrées !")
}
