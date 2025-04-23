package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"Back-end/db"
	"Back-end/db/dbmodels"

	"github.com/go-chi/chi"
)

func GetAllQuizzes(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, title, description, created_at, themes, id_user, id_game FROM quiz")
	if err != nil {
		http.Error(w, "Erreur lors de la récupération", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var quizzes []dbmodels.Quiz
	for rows.Next() {
		var q dbmodels.Quiz
		if err := rows.Scan(&q.ID, &q.Title, &q.Description, &q.CreatedAt, &q.Themes, &q.IdUser, &q.IdJeux); err != nil {
			http.Error(w, "Erreur de lecture des données", http.StatusInternalServerError)
			return
		}
		quizzes = append(quizzes, q)
	}
	json.NewEncoder(w).Encode(quizzes)
}

func GetQuizByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	row := db.DB.QueryRow("SELECT id, title, description, created_at, themes, id_user, id_game FROM quiz WHERE id = $1", id)
	var q dbmodels.Quiz
	if err := row.Scan(&q.ID, &q.Title, &q.Description, &q.CreatedAt, &q.Themes, &q.IdUser, &q.IdJeux); err != nil {
		http.Error(w, "Quiz introuvable", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(q)
}

func CreateQuiz(w http.ResponseWriter, r *http.Request) {
	var q dbmodels.Quiz
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		http.Error(w, "Format invalide", http.StatusBadRequest)
		return
	}
	_, err := db.DB.Exec(
		"INSERT INTO quiz (title, description, created_at, themes, id_user, id_game) VALUES ($1, $2, $3, $4, $5, $6)",
		q.Title, q.Description, q.CreatedAt, q.Themes, q.IdUser, q.IdJeux,
	)
	if err != nil {
		http.Error(w, "Erreur lors de la création", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(q)
}

func UpdateQuiz(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var q dbmodels.Quiz
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		http.Error(w, "Format invalide", http.StatusBadRequest)
		return
	}
	_, err := db.DB.Exec(
		"UPDATE quiz SET title=$1, description=$2, created_at=$3, themes=$4, id_user=$5, id_game=$6 WHERE id=$7",
		q.Title, q.Description, q.CreatedAt, q.Themes, q.IdUser, q.IdJeux, id,
	)
	if err != nil {
		http.Error(w, "Erreur lors de la mise à jour", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(q)
}

func DeleteQuiz(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idParsed, _ := strconv.Atoi(id)
	_, err := db.DB.Exec("DELETE FROM quiz WHERE id = $1", idParsed)
	if err != nil {
		http.Error(w, "Erreur lors de la suppression", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
