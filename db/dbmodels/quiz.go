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
	IdJeux      uint      `json:"id_game"`

	Questions []Question `json:"questions"`
}
