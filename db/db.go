package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/supabase-community/supabase-go"
)

var (
	DB       *sql.DB
	Supabase *supabase.Client
)

// InitDB initialise la connexion à la base de données (PostgreSQL)
func InitDB() {
	dbUrl := os.Getenv("SUPABASE_DB_URL")
	if dbUrl == "" {
		log.Fatal("SUPABASE_DB_URL non définie dans .env")
	}

	var err error
	DB, err = sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Erreur de connexion à la base de données : ", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal("Échec de la connexion à la base de données : ", err)
	}

	log.Println("✅ Connexion à la base de données réussie !")
}

// InitSupabase initialise le client Supabase
func InitSupabase() {
	var err error
	Supabase, err = supabase.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"), &supabase.ClientOptions{})
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation de Supabase : %v", err)
	}

	log.Println("✅ Client Supabase initialisé !")
}
