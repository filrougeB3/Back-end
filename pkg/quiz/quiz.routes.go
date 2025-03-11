package quiz

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Définir les routes pour les quiz
func RegisterRoutes(r *gin.Engine) {
	quizGroup := r.Group("/quiz")
	{
		quizGroup.GET("/", GetAllQuizzes)
		quizGroup.GET("/:id", GetQuizByID)
		quizGroup.POST("/create", CreateQuiz)
		quizGroup.PUT("/:id", UpdateQuiz)
		quizGroup.DELETE("/:id", DeleteQuiz)
	}
	fmt.Println("✅ Routes quiz enregistrées !") // Ajoute ce log pour confirmer
}
