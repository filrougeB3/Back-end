package propositions

import (
	"encoding/json"
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
	return strconv.Atoi(idStr)
}

func GetAllPropositions(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`SELECT id, value, is_correct, id_question FROM propositions`)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erreur récupération des propositions")
		return
	}
	defer rows.Close()

	var propositions []dbmodels.Proposition
	for rows.Next() {
		var p dbmodels.Proposition
		if err := rows.Scan(&p.ID, &p.Value, &p.IsCorrect, &p.IdQuestion); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Erreur lecture des propositions")
			return
		}
		propositions = append(propositions, p)
	}

	json.NewEncoder(w).Encode(propositions)
}

func GetPropositionByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDQueryParam(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID manquant ou invalide")
		return
	}

	var p dbmodels.Proposition
	err = db.DB.QueryRow(`SELECT id, value, is_correct, id_question FROM propositions WHERE id = $1`, id).
		Scan(&p.ID, &p.Value, &p.IsCorrect, &p.IdQuestion)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Proposition non trouvée")
		return
	}

	json.NewEncoder(w).Encode(p)
}

// POST /proposition/create
func CreateProposition(w http.ResponseWriter, r *http.Request) {
	var p dbmodels.Proposition
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

// PUT /proposition/update?id=1
func UpdateProposition(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDQueryParam(r)
	if err != nil || id == 0 {
		respondWithError(w, http.StatusBadRequest, "ID manquant ou invalide")
		return
	}

	var p dbmodels.Proposition
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

// DELETE /proposition/delete?id=1
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
