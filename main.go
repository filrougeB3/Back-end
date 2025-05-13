package main

import (
	"log"
	"net/http"

	"Back-end/db"
	"Back-end/pkg/auth"
	"Back-end/pkg/quiz"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// Initialiser la base de données et Supabase
	db.InitDB()
	db.InitSupabase()

	// Créer le routeur Chi
	router := chi.NewRouter()
	// Configuration CORS
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Autorise toutes les origines
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	router.Use(cors.Handler)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Enregistrement des routes
	auth.RegisterAuthRoutes(router)
	quiz.RegisterQuizRoutes(router)

	// Lancement du serveur
	log.Println("🚀 Le serveur tourne sur le port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
