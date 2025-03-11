package quiz

import (
	"fmt"
	"my-quiz-app/internal/service" // Assure-toi que le service existe dans le bon dossier
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetQuizzes(c *gin.Context) {
	// Appeler le service pour récupérer les quiz
	quizzes, err := service.GetAllQuizzes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la récupération des quiz"})
		return
	}
	c.JSON(http.StatusOK, quizzes)
}

func GetQuizDetails(c *gin.Context) {
	id := c.Param("id")
	quiz, err := service.GetQuizByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Quiz avec l'ID %s non trouvé", id)})
		return
	}
	c.JSON(http.StatusOK, quiz)
}

func CreateQuiz(c *gin.Context) {
	var quiz service.Quiz
	if err := c.ShouldBindJSON(&quiz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Données invalides"})
		return
	}

	// Appeler le service pour créer un quiz
	createdQuiz, err := service.CreateQuiz(quiz)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création du quiz"})
		return
	}

	c.JSON(http.StatusCreated, createdQuiz)
}

func UpdateQuiz(c *gin.Context) {
	id := c.Param("id")
	var quiz service.Quiz

	// Vérification de l'existence du quiz
	err := service.GetQuizByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Quiz avec l'ID %s non trouvé", id)})
		return
	}

	// Mise à jour des informations du quiz
	if err := c.ShouldBindJSON(&quiz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Données invalides"})
		return
	}

	// Appeler le service pour mettre à jour le quiz
	updatedQuiz, err := service.UpdateQuiz(id, quiz)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la mise à jour du quiz"})
		return
	}

	c.JSON(http.StatusOK, updatedQuiz)
}

func DeleteQuiz(c *gin.Context) {
	id := c.Param("id")

	// Vérifier si le quiz existe avant de le supprimer
	quiz, err := service.GetQuizByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Quiz avec l'ID %s non trouvé", id)})
		return
	}

	// Appeler le service pour supprimer le quiz
	err = service.DeleteQuiz(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la suppression du quiz"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Quiz avec l'ID %s supprimé avec succès", quiz.ID)})
}
