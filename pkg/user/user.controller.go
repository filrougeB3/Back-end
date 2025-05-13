package user

import (
	"Back-end/db"
	"Back-end/db/dbmodels"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, iduser, pseudo, email, country, profile_picture_url FROM users"
	rows, err := db.DB.Query(query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de la récupération des utilisateurs : %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []dbmodels.User

	for rows.Next() {
		var user dbmodels.User
		err := rows.Scan(
			&user.ID,
			&user.IDUser,
			&user.Pseudo,
			&user.Email,
			&user.Country,
			&user.ProfilePictureURL,
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("Erreur lors de la lecture d'un utilisateur : %v", err), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de l'itération des utilisateurs : %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// Récupérer le paramètre iduser dans l'URL
	idUser := chi.URLParam(r, "iduser")

	// Interroger la base de données pour récupérer l'utilisateur
	row := db.DB.QueryRow("SELECT id, iduser, pseudo, email, country, profile_picture_url FROM users WHERE iduser = $1", idUser)
	var user dbmodels.User
	if err := row.Scan(&user.ID, &user.IDUser, &user.Pseudo, &user.Email, &user.Country, &user.ProfilePictureURL); err != nil {
		http.Error(w, "Utilisateur introuvable", http.StatusNotFound)
		return
	}

	// Répondre avec les données de l'utilisateur
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID de l'utilisateur depuis l'URL
	idUser := chi.URLParam(r, "iduser")

	// Décoder le corps de la requête
	var updatedUser dbmodels.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Format de données invalide", http.StatusBadRequest)
		return
	}

	// Récupérer les données actuelles de l'utilisateur
	var currentUser dbmodels.User
	err := db.DB.QueryRow(
		"SELECT id, iduser, pseudo, email, country, profile_picture_url FROM users WHERE iduser = $1",
		idUser,
	).Scan(
		&currentUser.ID,
		&currentUser.IDUser,
		&currentUser.Pseudo,
		&currentUser.Email,
		&currentUser.Country,
		&currentUser.ProfilePictureURL,
	)

	if err != nil {
		http.Error(w, "Utilisateur non trouvé", http.StatusNotFound)
		return
	}

	// Préparer les valeurs à mettre à jour
	// Si un champ est vide dans le JSON, on garde la valeur actuelle
	pseudo := currentUser.Pseudo
	if updatedUser.Pseudo != nil && *updatedUser.Pseudo != "" {
		pseudo = updatedUser.Pseudo
	}

	email := currentUser.Email
	if updatedUser.Email != "" {
		email = updatedUser.Email
	}

	country := currentUser.Country
	if updatedUser.Country != nil && *updatedUser.Country != "" {
		country = updatedUser.Country
	}

	// Mettre à jour l'utilisateur
	query := `
		UPDATE users 
		SET pseudo = $1, email = $2, country = $3
		WHERE iduser = $4
		RETURNING id, iduser, pseudo, email, country, profile_picture_url
	`

	var user dbmodels.User
	err = db.DB.QueryRow(
		query,
		pseudo,
		email,
		country,
		idUser,
	).Scan(
		&user.ID,
		&user.IDUser,
		&user.Pseudo,
		&user.Email,
		&user.Country,
		&user.ProfilePictureURL,
	)

	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de la mise à jour de l'utilisateur : %v", err), http.StatusInternalServerError)
		return
	}

	// Répondre avec l'utilisateur mis à jour
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
