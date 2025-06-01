package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"Back-end/db"
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
	log.Println("➡️ URL appelée :", r.URL.Path)

	idStr := r.URL.Query().Get("id")
	log.Println("🧪 Paramètre ID (query) :", idStr)

	if idStr == "" {
		http.Error(w, "ID manquant dans l'URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	row := db.DB.QueryRow("SELECT id, title, description, created_at, themes, id_user, id_game FROM quiz WHERE id = $1", id)
	var q dbmodels.Quiz
	err = row.Scan(&q.ID, &q.Title, &q.Description, &q.CreatedAt, &q.Themes, &q.IdUser, &q.IdJeux)
	if err != nil {
		http.Error(w, "Quiz introuvable", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(q)
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
	idStr := r.URL.Query().Get("id")
	log.Println("🧪 Paramètre ID (query) :", idStr)

	if idStr == "" {
		http.Error(w, "ID manquant dans l'URL", http.StatusBadRequest)
		return
	}

	idParsed, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	var q dbmodels.Quiz
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		http.Error(w, "Format invalide", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec(
		"UPDATE quiz SET title=$1, description=$2, created_at=$3, themes=$4, id_user=$5, id_game=$6 WHERE id=$7",
		q.Title, q.Description, q.CreatedAt, q.Themes, q.IdUser, q.IdJeux, idParsed,
	)
	if err != nil {
		http.Error(w, "Erreur lors de la mise à jour", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(q)
}

func DeleteQuiz(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	log.Println("🧪 Paramètre ID (query) :", idStr)

	if idStr == "" {
		http.Error(w, "ID manquant dans l'URL", http.StatusBadRequest)
		return
	}

	idParsed, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	var exists bool
	err = db.DB.QueryRow("SELECT exists(SELECT 1 FROM quiz WHERE id = $1)", idParsed).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "Quiz introuvable", http.StatusNotFound)
		return
	}

	_, err = db.DB.Exec("DELETE FROM quiz WHERE id = $1", idParsed)
	if err != nil {
		http.Error(w, "Erreur lors de la suppression", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
