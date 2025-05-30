package dbmodels

type Question struct {
	ID            int    `json:"id"`
	Intitule      string `json:"intitulé"`
	IdQuiz        int    `json:"id_quiz"`
	IdType        int    `json:"id_type"`
	IdProposition int    `json:"id_proposition"`

	TypeQuestion *TypeQuestion `json:"type_question,omitempty"` // à remplir manuellement si besoin
	Propositions []Proposition `json:"propositions,omitempty"`  // idem
}
