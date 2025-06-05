package question

import (
	"Back-end/pkg/propositions"
)

// Modèle de la question
type Question struct {
	ID           int                        `json:"id"`
	Title        string                     `json:"title"`
	IdQuiz       int                        `json:"id_quiz"`
	IdType       int                        `json:"id_type"`
	Propositions []propositions.Proposition `json:"propositions,omitempty"`

	TypeQuestion *TypeQuestion `json:"type_question,omitempty"`
}

// Modèle du type de question
type TypeQuestion struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
