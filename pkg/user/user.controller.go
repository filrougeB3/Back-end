package user

import (
	"Back-end/db"
	"Back-end/pkg/security"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetUser godoc
// @Summary Récupérer les informations de l'utilisateur
// @Description Récupère les informations de l'utilisateur connecté
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} User
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(security.UserIDKey).(string)

	query := "SELECT iduser, email, pseudo, country, profile_picture_url FROM users WHERE iduser = $1"
	var user User
	err := db.DB.QueryRow(query, userID).Scan(&user.IDUser, &user.Email, &user.Pseudo, &user.Country, &user.Profile_picture_URL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de la récupération de l'utilisateur : %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// UpdateUser godoc
// @Summary Mettre à jour les informations de l'utilisateur
// @Description Met à jour les informations de l'utilisateur connecté
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body User true "Informations de l'utilisateur"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(security.UserIDKey).(string)

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors du décodage de l'utilisateur : %v", err), http.StatusBadRequest)
		return
	}

	query := "UPDATE users SET email = $1, pseudo = $2, country = $3, profile_picture_url = $4 WHERE iduser = $5"
	_, err = db.DB.Exec(query, user.Email, user.Pseudo, user.Country, user.Profile_picture_URL, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de la mise à jour de l'utilisateur : %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Utilisateur mis à jour"})
}
