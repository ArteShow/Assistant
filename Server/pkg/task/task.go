package task

import (
	"database/sql"
	"log"
	"os"
)

type Task struct {
	ID          int64  `json:"id"`
	Titel       string `json:"titel"`
	Status      string `json:"status"`
	Description string `json:"description"`
	UserID      int64  `json:"user_id"`
}

func NewTask(titel, description, status string) *Task {
	return &Task{
		Titel:       titel,
		Description: description,
		Status:      status,
	}
}

func SaveTask(task *Task, db *sql.DB) error {
	logFile, err := os.OpenFile("Server/log/task.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Printf("Inserting task '%s' for user ID %d into the database", task.Titel, task.UserID)

	query := `INSERT INTO tasks (titel, description, status, user_id) VALUES (?, ?, ?, ?)`
	_, err = db.Exec(query, task.Titel, task.Description, task.Status, task.UserID)
	if err != nil {
		log.Println("Error inserting task into database:", err)
	}
	return err
}

func GetAllUsersTasks(db *sql.DB, userID int64) ([]*Task, error) {
	log.Printf("Getting all tasks for user ID %d from the database", userID)

	query := `SELECT id, titel, description, status, user_id FROM tasks WHERE user_id = ?`
	rows, err := db.Query(query, userID)
	if err != nil {
		log.Println("Error getting tasks from database:", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Titel, &task.Description, &task.Status, &task.UserID)
		if err != nil {
			log.Println("Error scanning task:", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func GetUsersTaskByID(db *sql.DB, userID, taskID int64) (*Task, error) {
	log.Printf("Getting task with ID %d for user ID %d", taskID, userID)

	query := `SELECT id, titel, description, status, user_id FROM tasks WHERE user_id = ? AND id = ?`
	row := db.QueryRow(query, userID, taskID)

	task := &Task{}
	err := row.Scan(&task.ID, &task.Titel, &task.Description, &task.Status, &task.UserID)
	if err != nil {
		log.Println("Error scanning task:", err)
		return nil, err
	}

	return task, nil
}

func DeleteUsersTaskByID(db *sql.DB, userID, taskID int64) (bool, error) {
	log.Printf("Deleting task with ID %d for user ID %d", taskID, userID)

	query := `DELETE FROM tasks WHERE user_id = ? AND id = ?`
	result, err := db.Exec(query, userID, taskID)
	if err != nil {
		log.Println("Error deleting task:", err)
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected:", err)
		return false, err
	}

	return rowsAffected > 0, nil
}
