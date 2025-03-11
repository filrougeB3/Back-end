package quiz

import (
	"net/http"
	"strconv"

	database "Back-end/db/dbmodels"

	"github.com/gin-gonic/gin"
)

// Obtenir tous les quiz
func GetAllQuizzes(c *gin.Context) {
	quizzes, err := database.GetAllQuizzes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de récupérer les quiz"})
		return
	}
	c.JSON(http.StatusOK, quizzes)
}

// Obtenir un quiz par ID
func GetQuizByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	quiz, err := database.GetQuizByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur récupération quiz"})
		return
	}
	c.JSON(http.StatusOK, quiz)
}

// Créer un quiz
func CreateQuiz(c *gin.Context) {
	var quiz database.Quiz
	if err := c.ShouldBindJSON(&quiz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON invalide"})
		return
	}

	if err := database.CreateQuiz(quiz); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du quiz"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Quiz créé avec succès"})
}

// Modifier un quiz
func UpdateQuiz(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var quiz database.Quiz
	if err := c.ShouldBindJSON(&quiz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON invalide"})
		return
	}

	if err := database.UpdateQuiz(id, quiz); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du quiz"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Quiz mis à jour avec succès"})
}

// Supprimer un quiz
func DeleteQuiz(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	if err := database.DeleteQuiz(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression du quiz"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Quiz supprimé avec succès"})
}
