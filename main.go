package main

import (
	"Back-end/db"
	"Back-end/db/dbmodels"
	"Back-end/pkg/proposition"
	"Back-end/pkg/question"
	"Back-end/pkg/quiz"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Initialise la connexion à la base de données
	db.InitDB()

	// Utiliser GORM avec la même DATABASE_URL
	databaseUrl := db.GetDatabaseURL()
	var err error
	gormDB, err = gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})
	if err != nil {
		panic("Échec de la connexion GORM : " + err.Error())
	}

	dbmodels.Migrate(gormDB)
	db.SetGormDB(gormDB)

	r := gin.Default()

	quiz.RegisterQuizRoutes(r)
	question.RegisterQuestionRoutes(r)
	proposition.RegisterPropositionRoutes(r)

	// Lancer le serveur sur le port 8080
	fmt.Println("🚀 Serveur lancé sur http://localhost:8080")
	r.Run(":8080")
}
