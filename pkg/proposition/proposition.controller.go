package proposition

import (
	"encoding/json"
	"net/http"
	"strconv"

	"Back-end/db"
)

// GetAllPropositions godoc
// @Summary Récupérer toutes les propositions
// @Description Récupère la liste de toutes les propositions disponibles
// @Tags proposition
// @Accept json
// @Produce json
// @Success 200 {array} Proposition
// @Failure 500
// @Router /proposition/all [get]
func GetAllPropositions(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`SELECT id, value, is_correct, id_question FROM propositions`)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erreur récupération des propositions")
		return
	}
	defer rows.Close()

	var propositions []Proposition
	for rows.Next() {
		var p Proposition
		if err := rows.Scan(&p.ID, &p.Value, &p.IsCorrect, &p.IdQuestion); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Erreur lecture des propositions")
			return
		}
		propositions = append(propositions, p)
	}

	json.NewEncoder(w).Encode(propositions)
}

// GetPropositionByID godoc
// @Summary Récupérer une proposition par son ID
// @Description Récupère les détails d'une proposition spécifique
// @Tags proposition
// @Accept json
// @Produce json
// @Param id query int true "ID de la proposition"
// @Success 200 {object} Proposition
// @Failure 400
// @Failure 404
// @Router /proposition/get [get]
func GetPropositionByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDQueryParam(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID manquant ou invalide")
		return
	}

	var p Proposition
	err = db.DB.QueryRow(`SELECT id, value, is_correct, id_question FROM propositions WHERE id = $1`, id).
		Scan(&p.ID, &p.Value, &p.IsCorrect, &p.IdQuestion)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Proposition non trouvée")
		return
	}

	json.NewEncoder(w).Encode(p)
}

// CreateProposition godoc
// @Summary Créer une nouvelle proposition
// @Description Crée une nouvelle proposition
// @Tags proposition
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param proposition body Proposition true "Informations de la proposition"
// @Success 201 {object} Proposition
// @Failure 400
// @Failure 500
// @Router /proposition/create [post]
func CreateProposition(w http.ResponseWriter, r *http.Request) {
	var p Proposition
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Format invalide")
		return
	}

	err := db.DB.QueryRow(
		`INSERT INTO propositions (value, is_correct, id_question) VALUES ($1, $2, $3) RETURNING id`,
		p.Value, p.IsCorrect, p.IdQuestion,
	).Scan(&p.ID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erreur création proposition")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// UpdateProposition godoc
// @Summary Mettre à jour une proposition
// @Description Met à jour les informations d'une proposition existante
// @Tags proposition
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query int true "ID de la proposition"
// @Param proposition body Proposition true "Nouvelles informations de la proposition"
// @Success 200 {object} Proposition
// @Failure 400
// @Failure 500
// @Router /proposition/update [put]
func UpdateProposition(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDQueryParam(r)
	if err != nil || id == 0 {
		respondWithError(w, http.StatusBadRequest, "ID manquant ou invalide")
		return
	}

	var p Proposition
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Format invalide")
		return
	}

	_, err = db.DB.Exec(
		`UPDATE propositions SET value=$1, is_correct=$2, id_question=$3 WHERE id=$4`,
		p.Value, p.IsCorrect, p.IdQuestion, id,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erreur mise à jour proposition")
		return
	}

	p.ID = id
	json.NewEncoder(w).Encode(p)
}

// DeleteProposition godoc
// @Summary Supprimer une proposition
// @Description Supprime une proposition
// @Tags proposition
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query int true "ID de la proposition"
// @Success 204 "No Content"
// @Failure 400
// @Failure 500
// @Router /proposition/delete [delete]
func DeleteProposition(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDQueryParam(r)
	if err != nil || id == 0 {
		respondWithError(w, http.StatusBadRequest, "ID invalide")
		return
	}

	_, err = db.DB.Exec("DELETE FROM propositions WHERE id = $1", id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erreur suppression proposition")
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
	return strconv.Atoi(idStr)
}
