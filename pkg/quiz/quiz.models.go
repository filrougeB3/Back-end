package quiz

import (
	"Back-end/pkg/question"
	"time"
)

// Modèle du quiz
type Quiz struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Themes      string    `json:"themes"`
	IdUser      string    `json:"id_user"`
	IdJeux      int       `json:"id_game"`
	Pseudo      string    `json:"pseudo"` // Utilisé seulement pour les GET

	Questions []question.Question `json:"questions,omitempty"` // optionnel
}
