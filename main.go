package main

import (
	"Back-end/db"
	"Back-end/pkg/auth"
	"Back-end/pkg/quiz"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Charger les variables d'environnement depuis le fichier .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("❌ Erreur de chargement du fichier .env")
	}

	// Initialiser la base de données et Supabase
	db.InitDB()
	db.InitSupabase()

	// Créer le routeur Chi
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Enregistrement des routes
	auth.RegisterAuthRoutes(router)
	quiz.RegisterQuizRoutes(router)

	// Lancement du serveur
	log.Println("🚀 Le serveur tourne sur le port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
