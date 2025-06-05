package user

// Modèle de l'utilisateur
type User struct {
	IDUser              string  `json:"iduser"`
	Email               string  `json:"email"`
	Pseudo              string  `json:"pseudo"`
	Country             *string `json:"country"`
	Profile_picture_URL *string `json:"profile_picture_url"`
}
