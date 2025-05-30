package quiz

import (
	"Back-end/pkg/quiz/controller"
	"log"

	"github.com/go-chi/chi/v5"
)

func RegisterQuizRoutes(router chi.Router) {
	// Enregistrer les routes pour /quiz
	router.Route("/quiz", func(r chi.Router) {
		r.Get("/", controller.GetAllQuizzes)
		r.Post("/create", controller.CreateQuiz)

		// Routes utilisant {id} pour gérer un quiz spécifique
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", controller.GetQuizByID)
			r.Put("/", controller.UpdateQuiz)
			r.Delete("/", controller.DeleteQuiz)
		})
	})

	log.Println("✅ Routes quiz bien enregistrées !")
}
