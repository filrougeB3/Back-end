package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"Back-end/db/dbmodels"
	"Back-end/pkg/quiz/service"

	"github.com/go-chi/chi"
)

func parseIDFromRequest(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")
	log.Println("DEBUG path:", r.URL.Path)
	log.Println("DEBUG chi param id:", chi.URLParam(r, "id"))
	if idStr == "" {
		return 0, errors.New("ID manquant")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("ID invalide")
	}
	return id, nil
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func GetAllQuizzes(w http.ResponseWriter, r *http.Request) {
	quizzes, err := service.GetAllQuizzes()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erreur lors de la récupération")
		return
	}
	json.NewEncoder(w).Encode(quizzes)
}

func GetQuizByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromRequest(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	quiz, err := service.GetQuizWithQuestionsByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Quiz introuvable")
		return
	}
	json.NewEncoder(w).Encode(quiz)
}

func CreateQuiz(w http.ResponseWriter, r *http.Request) {
	var q dbmodels.Quiz
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		respondWithError(w, http.StatusBadRequest, "Format invalide")
		return
	}
	if err := service.ValidateQuiz(q); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := service.CreateQuiz(q); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erreur lors de la création")
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(q)
}

func UpdateQuiz(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromRequest(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	var q dbmodels.Quiz
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		respondWithError(w, http.StatusBadRequest, "Format invalide")
		return
	}
	q.ID = id
	if err := service.ValidateQuiz(q); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := service.UpdateQuiz(q); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erreur lors de la mise à jour")
		return
	}
	json.NewEncoder(w).Encode(q)
}

func DeleteQuiz(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromRequest(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := service.DeleteQuiz(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
