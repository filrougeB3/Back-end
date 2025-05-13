package dbmodels

type User struct {
	ID                int     `json:"id"`
	IDUser            string  `json:"iduser"`
	Pseudo            *string `json:"pseudo"`
	Email             string  `json:"email"`
	Country           *string `json:"country"`
	ProfilePictureURL *string `json:"profile_picture_url"`
}
