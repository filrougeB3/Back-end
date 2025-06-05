package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"Back-end/db"

	"github.com/supabase-community/gotrue-go/types"
)

// CreateUser godoc
// @Summary Créer un nouvel utilisateur
// @Description Crée un nouvel utilisateur avec email, pseudo et mot de passe
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "Informations de l'utilisateur"
// @Success 201 {object} RegisterResponse
// @Failure 400
// @Failure 500
// @Router /auth/register [post]
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

// LoginUser godoc
// @Summary Connexion utilisateur
// @Description Connecte un utilisateur avec email et mot de passe
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Identifiants de connexion"
// @Success 200 {object} LoginResponse
// @Failure 400
// @Failure 401
// @Router /auth/login [post]
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
