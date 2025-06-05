package question

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"Back-end/db"
)

// GetAllQuestions godoc
// @Summary Récupérer toutes les questions
// @Description Récupère la liste de toutes les questions disponibles
// @Tags question
// @Accept json
// @Produce json
// @Success 200 {array} Question
// @Failure 500
// @Router /question/all [get]
func GetAllQuestions(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`SELECT id, title, id_quiz, id_type FROM questions`)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erreur récupération des questions")
		return
	}
	defer rows.Close()

	var questions []Question
	for rows.Next() {
		var q Question
		if err := rows.Scan(&q.ID, &q.Title, &q.IdQuiz, &q.IdType); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Erreur lecture des questions")
			return
		}
		questions = append(questions, q)
	}

	json.NewEncoder(w).Encode(questions)
}

// GetQuestionByID godoc
// @Summary Récupérer une question par son ID
// @Description Récupère les détails d'une question spécifique
// @Tags question
// @Accept json
// @Produce json
// @Param id query int true "ID de la question"
// @Success 200 {object} Question
// @Failure 400
// @Failure 404
// @Router /question/get [get]
func GetQuestionByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDQueryParam(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID manquant ou invalide")
		return
	}

	var q Question
	err = db.DB.QueryRow(`SELECT id, title, id_quiz, id_type FROM questions WHERE id = $1`, id).
		Scan(&q.ID, &q.Title, &q.IdQuiz, &q.IdType)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Question non trouvée")
		return
	}

	json.NewEncoder(w).Encode(q)
}

// CreateQuestion godoc
// @Summary Créer une nouvelle question
// @Description Crée une nouvelle question
// @Tags question
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param question body Question true "Informations de la question"
// @Success 201 {object} Question
// @Failure 400
// @Failure 500
// @Router /question/create [post]
func CreateQuestion(w http.ResponseWriter, r *http.Request) {
	var q Question
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

// UpdateQuestion godoc
// @Summary Mettre à jour une question
// @Description Met à jour les informations d'une question existante
// @Tags question
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query int true "ID de la question"
// @Param question body Question true "Nouvelles informations de la question"
// @Success 200 {object} Question
// @Failure 400
// @Failure 500
// @Router /question/update [put]
func UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDQueryParam(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID manquant ou invalide")
		return
	}

	var q Question
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

// DeleteQuestion godoc
// @Summary Supprimer une question
// @Description Supprime une question et toutes ses propositions associées
// @Tags question
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query int true "ID de la question"
// @Success 204 "No Content"
// @Failure 400
// @Failure 500
// @Router /question/delete [delete]
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
