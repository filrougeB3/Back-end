package main

import (
	"Back-end/db" // Remplace par le chemin réel de ton package db
	"fmt"
)

func main() {
	// Initialise la connexion à la base de données
	db.InitDB()

	// Si la connexion est réussie, afficher un message
	fmt.Println("💻 L'application Go est prête et la connexion à la base est établie.")
}
