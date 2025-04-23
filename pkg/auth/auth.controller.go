package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"Back-end/db"

	"github.com/supabase-community/gotrue-go/types"
)

// Crée un nouvel utilisateur
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Données de la requête invalides", http.StatusBadRequest)
		return
	}

	signupRequest := types.SignupRequest{
		Email:    req.Email,
		Password: req.Password,
	}
	authResp, err := db.Supabase.Auth.Signup(signupRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de la création de l'utilisateur : %v", err), http.StatusInternalServerError)
		return
	}

	uid := authResp.User.ID

	query := "INSERT INTO users (iduser, email, pseudo) VALUES ($1, $2, $3) RETURNING iduser"
	var iduser string
	err = db.DB.QueryRow(query, uid, req.Email, req.Pseudo).Scan(&iduser)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de l'insertion dans la base de données : %v", err), http.StatusInternalServerError)
		return
	}

	resp := RegisterResponse{
		IDUser:       iduser,
		Email:        req.Email,
		Pseudo:       req.Pseudo,
		AuthToken:    authResp.AccessToken,
		RefreshToken: authResp.RefreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// Connexion de l'utilisateur
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Données de la requête invalides", http.StatusBadRequest)
		return
	}

	authResp, err := db.Supabase.Auth.SignInWithEmailPassword(req.Email, req.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de la connexion de l'utilisateur : %v", err), http.StatusUnauthorized)
		return
	}

	resp := LoginResponse{
		Email:        req.Email,
		AuthToken:    authResp.AccessToken,
		RefreshToken: authResp.RefreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Déconnexion de l'utilisateur
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	// Renvoyer simplement le message de déconnexion
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Utilisateur déconnecté"})
}
