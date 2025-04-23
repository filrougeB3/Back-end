package dbmodels

type Question struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	Intitule      string `json:"intitulé"`
	IdQuiz        uint   `json:"id_quiz"`
	IdType        uint   `json:"id_type"`
	IdProposition uint   `json:"id_proposition"`

	TypeQuestion TypeQuestion  `gorm:"foreignKey:IdType" json:"type_question"` // Relation avec TypeQuestion
	Propositions []Proposition `gorm:"foreignKey:IdQuestion" json:"propositions"`
}
