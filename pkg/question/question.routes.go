package question

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterQuestionRoutes(router chi.Router) {
	router.Route("/question", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Get("/all", GetAllQuestions)
		r.Get("/get", GetQuestionByID)
		r.Post("/create", CreateQuestion)
		r.Put("/update", UpdateQuestion)
		r.Delete("/delete", DeleteQuestion)
	})
	log.Println("✅ Routes question enregistrées !")
}
