package quiz

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"Back-end/db"
	"Back-end/db/dbmodels"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func parseIDQueryParam(r *http.Request) (int, error) {
	idStr := r.URL.Query().Get("id")
	log.Println("🧪 Paramètre ID (query) :", idStr)

	if idStr == "" {
		return 0, errors.New("ID manquant dans l'URL")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("ID invalide")
	}

	return id, nil
}

// GET /quiz/all
func GetAllQuizzes(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, title, description, created_at, themes, id_user, id_game FROM quiz")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erreur lors de la récupération")
		return
	}
	defer rows.Close()

	var quizzes []dbmodels.Quiz
	for rows.Next() {
		var q dbmodels.Quiz
		if err := rows.Scan(&q.ID, &q.Title, &q.Description, &q.CreatedAt, &q.Themes, &q.IdUser, &q.IdJeux); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Erreur lors du scan")
			return
		}
		quizzes = append(quizzes, q)
	}

	json.NewEncoder(w).Encode(quizzes)
}

// GET /quiz?id=1
func GetQuizByID(w http.ResponseWriter, r *http.Request) {
	log.Println("➡️ URL appelée :", r.URL.Path)

	id, err := parseIDQueryParam(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

// POST /quiz/create
func CreateQuiz(w http.ResponseWriter, r *http.Request) {
	var q dbmodels.Quiz
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		respondWithError(w, http.StatusBadRequest, "Format invalide")
		return
	}

	err := db.DB.QueryRow(
		`INSERT INTO quiz (title, description, created_at, themes, id_user, id_game) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		q.Title, q.Description, q.CreatedAt, q.Themes, q.IdUser, q.IdJeux,
	).Scan(&q.ID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erreur lors de la création")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(q)
}

// PUT /quiz/update?id=1
func UpdateQuiz(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDQueryParam(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var q dbmodels.Quiz
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		http.Error(w, "Format invalide", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec(
		`UPDATE quiz 
		SET title=$1, description=$2, created_at=$3, themes=$4, id_user=$5, id_game=$6 
		WHERE id=$7`,
		q.Title, q.Description, q.CreatedAt, q.Themes, q.IdUser, q.IdJeux, id,
	)
	if err != nil {
		http.Error(w, "Erreur lors de la mise à jour", http.StatusInternalServerError)
		return
	}

	q.ID = id
	json.NewEncoder(w).Encode(q)
}

// DELETE /quiz/delete?id=1
func DeleteQuiz(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDQueryParam(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var exists bool
	err = db.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM quiz WHERE id = $1)", id).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "Quiz introuvable", http.StatusNotFound)
		return
	}

	_, err = db.DB.Exec("DELETE FROM quiz WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Erreur lors de la suppression", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
