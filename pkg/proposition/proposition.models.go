package proposition

type Proposition struct {
	ID         int    `json:"id"`
	Value      string `json:"value"`
	IsCorrect  bool   `json:"is_correct"`
	IdQuestion int    `json:"id_question"`
}
