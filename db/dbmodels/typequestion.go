package dbmodels

type TypeQuestion struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Valeur      string `json:"valeur"`
	Description string `json:"description"`

	// Pas besoin de id_question ici
}
