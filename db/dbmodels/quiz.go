package dbmodels

import (
	"time"
)

type Quiz struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Themes      string    `json:"themes"`
	IdUser      uint      `json:"id_user"`
	IdJeux      uint      `json:"id_jeux"`

	Questions []Question `gorm:"foreignKey:IdQuiz" json:"questions"` // Relation 1-n avec Question
}
