package service

import (
	"errors"
	"time"

	"Back-end/db"
	"Back-end/db/dbmodels"
)

func GetAllQuizzes() ([]dbmodels.Quiz, error) {
	rows, err := db.DB.Query("SELECT id, title, description, created_at, themes, id_user, id_game FROM quiz")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes []dbmodels.Quiz
	for rows.Next() {
		var q dbmodels.Quiz
		if err := rows.Scan(&q.ID, &q.Title, &q.Description, &q.CreatedAt, &q.Themes, &q.IdUser, &q.IdJeux); err != nil {
			return nil, err
		}
		quizzes = append(quizzes, q)
	}
	return quizzes, nil
}

func GetQuizWithQuestionsByID(id int) (dbmodels.Quiz, error) {
	var q dbmodels.Quiz
	err := db.DB.QueryRow("SELECT id, title, description, created_at, themes, id_user, id_game FROM quiz WHERE id = $1", id).
		Scan(&q.ID, &q.Title, &q.Description, &q.CreatedAt, &q.Themes, &q.IdUser, &q.IdJeux)
	if err != nil {
		return dbmodels.Quiz{}, err
	}

	questions, err := getQuestionsByQuizID(uint(id))
	if err == nil {
		q.Questions = questions
	}
	return q, nil
}

func CreateQuiz(q dbmodels.Quiz) error {
	_, err := db.DB.Exec(
		"INSERT INTO quiz (title, description, created_at, themes, id_user, id_game) VALUES ($1, $2, $3, $4, $5, $6)",
		q.Title, q.Description, time.Now(), q.Themes, q.IdUser, q.IdJeux,
	)
	return err
}

func UpdateQuiz(q dbmodels.Quiz) error {
	_, err := db.DB.Exec(
		"UPDATE quiz SET title=$1, description=$2, created_at=$3, themes=$4, id_user=$5, id_game=$6 WHERE id=$7",
		q.Title, q.Description, q.CreatedAt, q.Themes, q.IdUser, q.IdJeux, q.ID,
	)
	return err
}

func DeleteQuiz(id int) error {
	var exists bool
	err := db.DB.QueryRow("SELECT exists(SELECT 1 FROM quiz WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("quiz non trouvé")
	}

	_, err = db.DB.Exec("DELETE FROM quiz WHERE id = $1", id)
	return err
}

func ValidateQuiz(q dbmodels.Quiz) error {
	if q.Title == "" {
		return errors.New("le titre est obligatoire")
	}
	if q.Description == "" {
		return errors.New("la description est obligatoire")
	}
	if q.IdUser == "" {
		return errors.New("l'utilisateur est requis")
	}
	if q.IdJeux <= 0 {
		return errors.New("le jeu est invalide")
	}
	return nil
}

func getQuestionsByQuizID(quizID uint) ([]dbmodels.Question, error) {
	rows, err := db.DB.Query("SELECT id, intitule, id_quiz, id_type, id_proposition FROM question WHERE id_quiz = $1", quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []dbmodels.Question
	for rows.Next() {
		var q dbmodels.Question
		if err := rows.Scan(&q.ID, &q.Intitule, &q.IdQuiz, &q.IdType, &q.IdProposition); err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}
	return questions, nil
}
