package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var Conn *pgx.Conn

// InitDB initialise la connexion à la base de données
func InitDB() {
	// Charger les variables d'environnement depuis le fichier .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erreur lors du chargement du fichier .env")
	}

	// Récupérer l'URL de la base de données
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatal("DATABASE_URL non définie dans .env")
	}

	// Connexion à PostgreSQL
	conn, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		log.Fatalf("Impossible de se connecter à la base : %v", err)
	}

	fmt.Println("✅ Connexion à PostgreSQL réussie !")
	Conn = conn
}
