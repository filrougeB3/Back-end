package main

import (
	"Back-end/db"
	"Back-end/pkg/auth"
	"Back-end/pkg/proposition"
	"Back-end/pkg/question"
	"Back-end/pkg/quiz"
	"Back-end/pkg/user"
	"log"
	"net/http"

	_ "Back-end/docs" // Ceci est important pour charger la documentation Swagger

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title NoLedge API
// @version 1.0
// @description Backend de l'application NoLedge, une plateforme de quiz en ligne
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @x-extension-openapi {"servers":[{"url":"http://localhost:8080","description":"Serveur local"},{"url":"https://back-end-73xk.onrender.com","description":"Serveur de production"}]}
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

	// Swagger UI
	router.Get("/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// Enregistrement des routes
	auth.RegisterAuthRoutes(router)
	quiz.RegisterQuizRoutes(router)
	user.RegisterUserRoutes(router)
	question.RegisterQuestionRoutes(router)
	proposition.RegisterPropositionRoutes(router)

	// Lancement du serveur
	log.Println("🚀 Le serveur tourne sur le port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
