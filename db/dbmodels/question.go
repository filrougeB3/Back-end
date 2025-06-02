package dbmodels

type Question struct {
	ID           int           `json:"id"`
	Title        string        `json:"title"`
	IdQuiz       int           `json:"id_quiz"`
	IdType       int           `json:"id_type"`
	Propositions []Proposition `json:"propositions,omitempty"`

	TypeQuestion *TypeQuestion `json:"type_question,omitempty"`
}
