package task

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID          int64
	Titel       string
	Status      string
	Description string
	userID    int64
}

type User struct {
	UserID int64
	Tasks  []Task
}

func NewTask(titel, description, status string) *Task {
	return &Task{
		Titel:       titel,
		Description: description,
		Status:      status,
	}
}

func SaveTask(task *Task, userID int64, db *sql.DB) error {
	log.Printf("Inserting task '%s' for user ID %d into the database", task.Titel, userID)

	query := `INSERT INTO tasks (titel, description, status, user_id) VALUES (?, ?, ?, ?)`

	_, err := db.Exec(query, task.Titel, task.Description, task.Status, userID)
	if err != nil {
		log.Println("Error inserting task into database:", err)
		return err
	}

	return nil
}

func GetAllUsersTasks(db *sql.DB, userID int64) ([]*Task, error) {
	log.Printf("Getting all tasks for user ID %d from the database", userID)

	query := `SELECT * FROM tasks WHERE user_id = ?`
	rows, err := db.Query(query, userID)
	if err != nil {
		log.Println("Error getting tasks from database:", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Titel, &task.Description, &task.Status)
		if err != nil {
			log.Println("Error scanning task:", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func GetUsersTaskByID(userID int64, taskID int64, db *sql.DB) (*Task, error) {
    log.Printf("Getting task with ID %d for user ID %d", taskID, userID)

    query := `SELECT id, titel, description, status FROM tasks WHERE user_id = ? AND id = ?`
    rows, err := db.Query(query, userID, taskID)
    if err != nil {
        log.Println("Error getting task from database:", err)
        return nil, err
    }
    defer rows.Close()

    var task *Task
    if rows.Next() {
        task = &Task{}
        err := rows.Scan(&task.ID, &task.Titel, &task.Description, &task.Status)
        if err != nil {
            log.Println("Error scanning task:", err)
            return nil, err
        }
    } else {
        log.Printf("No task found with ID %d for user ID %d", taskID, userID)
        return nil, sql.ErrNoRows
    }

    return task, nil
}

func DeletUsersTaskByID(userID int64, taskID int64, db *sql.DB) (bool, error){
	log.Printf("Deleting task with ID %d for user ID %d", taskID, userID)

	query := `DELETE FROM tasks WHERE user_id = ? AND id = ?`
	result, err := db.Exec(query, userID, taskID)
	if err != nil {
		log.Println("Error deleting task from database:", err)
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return false, err
	}

	if rowsAffected == 0 {
		log.Printf("No task found with ID %d for user ID %d", taskID, userID)
		return false, nil
	}

	return true, nil
}