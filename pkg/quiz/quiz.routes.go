package quiz

import (
	"Back-end/pkg/quiz/controller"

	"github.com/gin-gonic/gin"
)

func RegisterQuizRoutes(r *gin.Engine) {
	quizGroup := r.Group("/quiz")
	{
		quizGroup.GET("", controller.GetAllQuizzes)
		quizGroup.GET(":id", controller.GetQuizByID)
		quizGroup.POST("/create", controller.CreateQuiz)
		quizGroup.PUT(":id", controller.UpdateQuiz)
		quizGroup.DELETE(":id", controller.DeleteQuiz)
	}
}
