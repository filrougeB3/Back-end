package quiz

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"Back-end/db"
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

// GetAllQuizzes godoc
// @Summary Récupérer tous les quiz
// @Description Récupère la liste de tous les quiz disponibles
// @Tags quiz
// @Accept json
// @Produce json
// @Success 200 {array} Quiz
// @Failure 500
// @Router /quiz/all [get]
func GetAllQuizzes(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`
		SELECT q.id, q.title, q.description, q.created_at, q.themes, q.id_user, q.id_game, u.pseudo 
		FROM quiz q
		LEFT JOIN users u ON q.id_user = u.iduser
	`)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erreur lors de la récupération")
		return
	}
	defer rows.Close()

	var quizzes []Quiz
	for rows.Next() {
		var q Quiz
		if err := rows.Scan(&q.ID, &q.Title, &q.Description, &q.CreatedAt, &q.Themes, &q.IdUser, &q.IdJeux, &q.Pseudo); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Erreur lors du scan")
			return
		}
		quizzes = append(quizzes, q)
	}

	json.NewEncoder(w).Encode(quizzes)
}

// GetQuizByID godoc
// @Summary Récupérer un quiz par son ID
// @Description Récupère les détails d'un quiz spécifique
// @Tags quiz
// @Accept json
// @Produce json
// @Param id query int true "ID du quiz"
// @Success 200 {object} Quiz
// @Failure 400
// @Failure 404
// @Router /quiz/byQuery [get]
func GetQuizByID(w http.ResponseWriter, r *http.Request) {
	log.Println("➡️ URL appelée :", r.URL.Path)

	id, err := parseIDQueryParam(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	row := db.DB.QueryRow(`
		SELECT q.id, q.title, q.description, q.created_at, q.themes, q.id_user, q.id_game, u.pseudo 
		FROM quiz q
		LEFT JOIN users u ON q.id_user = u.iduser
		WHERE q.id = $1
	`, id)
	var q Quiz
	err = row.Scan(&q.ID, &q.Title, &q.Description, &q.CreatedAt, &q.Themes, &q.IdUser, &q.IdJeux, &q.Pseudo)
	if err != nil {
		http.Error(w, "Quiz introuvable", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(q)
}

// CreateQuiz godoc
// @Summary Créer un nouveau quiz
// @Description Crée un nouveau quiz avec ses questions et propositions
// @Tags quiz
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param quiz body Quiz true "Informations du quiz"
// @Success 201 {object} Quiz
// @Failure 400
// @Failure 500
// @Router /quiz/create [post]
func CreateQuiz(w http.ResponseWriter, r *http.Request) {
	var q Quiz
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		respondWithError(w, http.StatusBadRequest, "Format invalide")
		return
	}

	// Démarrer transaction
	tx, err := db.DB.Begin()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erreur transaction")
		return
	}

	// Insert quiz et récupère ID
	err = tx.QueryRow(
		`INSERT INTO quiz (title, description, created_at, themes, id_user, id_game) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		q.Title, q.Description, q.CreatedAt, q.Themes, q.IdUser, q.IdJeux,
	).Scan(&q.ID)

	if err != nil {
		tx.Rollback()
		respondWithError(w, http.StatusInternalServerError, "Erreur lors de la création du quiz")
		return
	}

	// Insert questions et propositions si présentes
	for i, question := range q.Questions {
		var questionID int
		err = tx.QueryRow(
			`INSERT INTO questions (title, id_quiz, id_type) VALUES ($1, $2, $3) RETURNING id`,
			question.Title, q.ID, question.IdType,
		).Scan(&questionID)
		if err != nil {
			log.Printf("💥 Erreur SQL lors de l'insertion question : %+v", err)
			respondWithError(w, http.StatusInternalServerError, "Erreur lors de la création d'une question")
			tx.Rollback()
			return
		}

		// update l'ID dans la structure pour retour JSON
		q.Questions[i].ID = questionID

		// insert propositions pour la question
		for j, prop := range question.Propositions {
			var propID int
			err = tx.QueryRow(
				`INSERT INTO propositions (value, is_correct, id_question) VALUES ($1, $2, $3) RETURNING id`,
				prop.Value, prop.IsCorrect, questionID,
			).Scan(&propID)
			if err != nil {
				tx.Rollback()
				respondWithError(w, http.StatusInternalServerError, "Erreur lors de la création d'une proposition")
				return
			}

			q.Questions[i].Propositions[j].ID = propID
			q.Questions[i].Propositions[j].IdQuestion = questionID
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erreur lors du commit")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(q)
}

// UpdateQuiz godoc
// @Summary Mettre à jour un quiz
// @Description Met à jour les informations d'un quiz existant
// @Tags quiz
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query int true "ID du quiz"
// @Param quiz body Quiz true "Nouvelles informations du quiz"
// @Success 200 {object} Quiz
// @Failure 400
// @Failure 500
// @Router /quiz/byQuery [put]
func UpdateQuiz(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDQueryParam(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var q Quiz
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

// DeleteQuiz godoc
// @Summary Supprimer un quiz
// @Description Supprime un quiz et toutes ses questions associées
// @Tags quiz
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query int true "ID du quiz"
// @Success 204 "No Content"
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /quiz/byQuery [delete]
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
