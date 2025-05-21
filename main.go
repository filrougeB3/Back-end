package main

import (
	"log"
	"net/http"

	"Back-end/db"
	"Back-end/pkg/auth"
	"Back-end/pkg/quiz"
	"Back-end/pkg/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erreur lors du chargement du fichier .env :", err)
	}
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
	user.RegisterUserRoutes(router)

	// Lancement du serveur
	log.Println("🚀 Le serveur tourne sur le port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
