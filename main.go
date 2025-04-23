package main

import (
	"log"
	"net/http"

	"Back-end/db"
	"Back-end/pkg/auth"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv" // Importer godotenv
)

func main() {
	// Charger les variables d'environnement depuis le fichier .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erreur de chargement du fichier .env")
	}

	// Initialiser la connexion à la base de données
	db.InitDB()
	db.InitSupabase()

	// Créer un nouveau routeur
	router := mux.NewRouter()

	// Définir les routes d'authentification
	auth.AuthRoute(router)

	// Démarrer le serveur HTTP
	log.Println("🚀 Le serveur tourne sur le port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
