package dbmodels

type TypeQuestion struct {
	ID          uint   `json:"id"`
	Valeur      string `json:"valeur"`
	Description string `json:"description"`
}
