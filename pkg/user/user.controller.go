package user

import (
	"Back-end/db"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)

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

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)

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