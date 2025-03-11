package quiz

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes est une fonction qui enregistre toutes les routes liées aux quiz
func RegisterRoutes(r *gin.Engine) {
	quizGroup := r.Group("/quiz")
	{
		quizGroup.GET("/", GetQuizzes)
		quizGroup.GET("/:id", GetQuizDetails)
		quizGroup.POST("/create", CreateQuiz)
		quizGroup.PUT("/:id", UpdateQuiz)
		quizGroup.DELETE("/:id", DeleteQuiz)
	}
}
