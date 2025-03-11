package database

import (
	"context"
	"log"

	"Back-end/db"

	"github.com/jackc/pgx/v5"
)

// Structure du modèle Quiz
type Quiz struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Theme     string `json:"theme"`
	Questions string `json:"questions"`
	CreatedBy int    `json:"created_by"`
}

// Récupérer tous les quiz
func GetAllQuizzes() ([]Quiz, error) {
	rows, err := db.Conn.Query(context.Background(), "SELECT id, title, theme, questions, created_by FROM quiz")
	if err != nil {
		log.Printf("Erreur récupération quiz : %v", err)
		return nil, err
	}
	defer rows.Close()

	var quizzes []Quiz
	for rows.Next() {
		var q Quiz
		err := rows.Scan(&q.ID, &q.Title, &q.Theme, &q.Questions, &q.CreatedBy)
		if err != nil {
			log.Printf("Erreur lecture ligne : %v", err)
			continue
		}
		quizzes = append(quizzes, q)
	}

	return quizzes, nil
}

// Récupérer un quiz par ID
func GetQuizByID(id int) (Quiz, error) {
	var quiz Quiz
	err := db.Conn.QueryRow(context.Background(),
		"SELECT id, title, theme, questions, created_by FROM quiz WHERE id=$1", id).
		Scan(&quiz.ID, &quiz.Title, &quiz.Theme, &quiz.Questions, &quiz.CreatedBy)

	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("Aucun quiz trouvé avec ID %d", id)
			return quiz, nil
		}
		log.Printf("Erreur récupération quiz ID %d : %v", id, err)
		return quiz, err
	}

	return quiz, nil
}

// Créer un quiz
func CreateQuiz(quiz Quiz) error {
	_, err := db.Conn.Exec(context.Background(),
		"INSERT INTO quiz (title, theme, questions, created_by) VALUES ($1, $2, $3, $4)",
		quiz.Title, quiz.Theme, quiz.Questions, quiz.CreatedBy)

	if err != nil {
		log.Printf("Erreur création quiz : %v", err)
	}

	return err
}

// Modifier un quiz
func UpdateQuiz(id int, quiz Quiz) error {
	_, err := db.Conn.Exec(context.Background(),
		"UPDATE quiz SET title=$1, theme=$2, questions=$3, created_by=$4 WHERE id=$5",
		quiz.Title, quiz.Theme, quiz.Questions, quiz.CreatedBy, id)

	if err != nil {
		log.Printf("Erreur mise à jour quiz ID %d : %v", id, err)
	}

	return err
}

// Supprimer un quiz
func DeleteQuiz(id int) error {
	_, err := db.Conn.Exec(context.Background(),
		"DELETE FROM quiz WHERE id=$1", id)

	if err != nil {
		log.Printf("Erreur suppression quiz ID %d : %v", id, err)
	}

	return err
}
