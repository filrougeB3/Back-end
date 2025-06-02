package question

import (
	"encoding/json"
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
	if idStr == "" {
		return 0, nil
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func CreateQuestion(w http.ResponseWriter, r *http.Request) {
	var q dbmodels.Question
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		respondWithError(w, http.StatusBadRequest, "Format invalide")
		return
	}

	err := db.DB.QueryRow(
		`INSERT INTO questions (title, id_quiz, id_type) VALUES ($1, $2, $3) RETURNING id`,
		q.Title, q.IdQuiz, q.IdType,
	).Scan(&q.ID)
	if err != nil {
		log.Printf("💥 Erreur SQL lors de l'insertion question : %+v", err)
		respondWithError(w, http.StatusInternalServerError, "Erreur création question")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(q)
}

func UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDQueryParam(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID manquant ou invalide")
		return
	}

	var q dbmodels.Question
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		respondWithError(w, http.StatusBadRequest, "Format invalide")
		return
	}

	_, err = db.DB.Exec(
		`UPDATE questions SET title=$1, id_quiz=$2, id_type=$3 WHERE id=$4`,
		q.Title, q.IdQuiz, q.IdType, id,
	)
	if err != nil {
		log.Printf("💥 Erreur SQL lors de l'insertion question : %+v", err)
		respondWithError(w, http.StatusInternalServerError, "Erreur update question")
		return
	}

	q.ID = id
	json.NewEncoder(w).Encode(q)
}

func DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDQueryParam(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID invalide")
		return
	}

	// Supprimer d'abord les propositions liées
	_, err = db.DB.Exec("DELETE FROM propositions WHERE id_question = $1", id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erreur suppression propositions")
		return
	}

	// Ensuite la question
	_, err = db.DB.Exec("DELETE FROM questions WHERE id = $1", id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erreur suppression question")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
