package dbmodels

type Proposition struct {
	ID         int    `json:"id"`
	Valeur     string `json:"valeur"`
	IsCorrect  bool   `json:"is_correct"`
	IdQuestion int    `json:"id_question"`

	// Pas de liaison directe, à gérer manuellement si besoin
}
