package auth

// Requête de création de compte
type RegisterRequest struct {
	Email    string `json:"email"`
	Pseudo   string `json:"pseudo"`
	Password string `json:"password"`
}

// Réponse de la création de compte
type RegisterResponse struct {
	IDUser       string `json:"iduser"`
	Email        string `json:"email"`
	Pseudo       string `json:"pseudo"`
	AuthToken    string `json:"auth_token"`
	RefreshToken string `json:"refresh_token"`
}

// Requête de connexion
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Réponse de la connexion
type LoginResponse struct {
	Email        string `json:"email"`
	AuthToken    string `json:"auth_token"`
	RefreshToken string `json:"refresh_token"`
}
