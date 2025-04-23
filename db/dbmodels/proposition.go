package dbmodels

type Proposition struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Valeur     string `json:"valeur"`
	IsCorrect  bool   `json:"is_correct"`
	IdQuestion uint   `json:"id_question"`

	Question Question `gorm:"foreignKey:IdQuestion" json:"-"`
}
