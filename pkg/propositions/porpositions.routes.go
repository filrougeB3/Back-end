package propositions

import (
	"log"

	"github.com/go-chi/chi/v5"
)

func RegisterPropositionRoutes(router chi.Router) {
	router.Route("/proposition", func(r chi.Router) {
		r.Post("/create", CreateProposition)
		r.Put("/update", UpdateProposition)
		r.Delete("/delete", DeleteProposition)
	})
	log.Println("✅ Routes proposition enregistrées !")
}
