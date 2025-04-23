package controller

import (
	"net/http"
	"strconv"

	"Back-end/db"
	"Back-end/db/dbmodels"

	"github.com/gin-gonic/gin"
)

func GetAllQuizzes(c *gin.Context) {
	var quizzes []dbmodels.Quiz
	if err := db.GetGormDB().Preload("Questions.Propositions").Find(&quizzes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération"})
		return
	}
	c.JSON(http.StatusOK, quizzes)
}

func GetQuizByID(c *gin.Context) {
	id := c.Param("id")
	var quiz dbmodels.Quiz
	if err := db.GetGormDB().Preload("Questions.Propositions").First(&quiz, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz introuvable"})
		return
	}
	c.JSON(http.StatusOK, quiz)
}

func CreateQuiz(c *gin.Context) {
	var quiz dbmodels.Quiz
	if err := c.ShouldBindJSON(&quiz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format invalide"})
		return
	}
	if err := db.GetGormDB().Create(&quiz).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création"})
		return
	}
	c.JSON(http.StatusCreated, quiz)
}

func UpdateQuiz(c *gin.Context) {
	id := c.Param("id")
	var quiz dbmodels.Quiz
	if err := db.GetGormDB().First(&quiz, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz introuvable"})
		return
	}
	if err := c.ShouldBindJSON(&quiz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format invalide"})
		return
	}
	db.GetGormDB().Save(&quiz)
	c.JSON(http.StatusOK, quiz)
}

func DeleteQuiz(c *gin.Context) {
	id := c.Param("id")
	idParsed, _ := strconv.Atoi(id)
	db.GetGormDB().Delete(&dbmodels.Quiz{}, idParsed)
	c.Status(http.StatusNoContent)
}
