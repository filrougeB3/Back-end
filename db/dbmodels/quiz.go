package dbmodels

import (
	"time"
)

type Quiz struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Themes      string    `json:"themes"`
	IdUser      string    `json:"id_user"`
	IdJeux      int       `json:"id_game"`

	Questions []Question `json:"questions"`
}
