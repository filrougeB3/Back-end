package user

import (
	"github.com/go-chi/chi/v5"
)

func UserRoutes(router chi.Router) {
	router.Route("/user", func(r chi.Router) {
		r.Get("/", GetAllUsers)
		r.Get("/{iduser}", GetUser)               // Récupérer un utilisateur par ID
		r.Put("/updateUser/{iduser}", UpdateUser) // Mettre à jour un utilisateur
	})
}
