package authorization

import (
	"database/sql"
	"log"
)

type User struct {
	UserID   int64
	Username string
	Password string
	Email    string
}

func NewUser(username, password, email string) *User {
	return &User{
		Username: username,
		Password: password,
		Email:    email,
	}
}

func CheckForUser(userID int64, db *sql.DB) (bool, error) {
	log.Printf("Checking if user with ID %d exists in the database", userID)

	query := `SELECT COUNT(*) FROM users WHERE user_id = ?`
	row := db.QueryRow(query, userID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		log.Println("Error checking user in database:", err)
		return false, err
	}
	if count > 0 || count == 1 {
		log.Printf("User with ID %d exists in the database", userID)
		return true, nil
	}
	return false, nil
}

func SaveUser(user *User, db *sql.DB) error {
	log.Printf("Inserting user '%s' into the database", user.Username)

	query := `INSERT INTO users (username, password) VALUES (?, ?)`

	_, err := db.Exec(query, user.Username, user.Password, user.Email)
	if err != nil {
		log.Println("Error inserting user into database:", err)
		return err
	}

	return nil
}

func DeletUser(userID int64, db *sql.DB) error {
	log.Printf("Deleting user with ID %d from the database", userID)

	query := `DELETE FROM users WHERE user_id = ?`

	_, err := db.Exec(query, userID)
	if err != nil {
		log.Println("Error deleting user from database:", err)
		return err
	}

	return nil
}

func GetUserById(userId int64, db *sql.DB) (*User, error) {
	log.Printf("Getting user with ID %d from the database", userId)

	query := `SELECT * FROM users WHERE user_id = ?`
	row := db.QueryRow(query, userId)

	user := &User{}
	err := row.Scan(&user.UserID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		log.Println("Error getting user from database:", err)
		return nil, err
	}

	return user, nil
}

func ChangeUsersPassword(userID int64, newPassword string, db *sql.DB) error {
	log.Printf("Changing password for user with ID %d", userID)

	query := `UPDATE users SET password = ? WHERE user_id = ?`

	_, err := db.Exec(query, newPassword, userID)
	if err != nil {
		log.Println("Error changing user password in database:", err)
		return err
	}

	return nil
}

func ChangeUsersUsername(userID int64, newUsername string, db *sql.DB) error {
	log.Printf("Changing username for user with ID %d", userID)

	query := `UPDATE users SET username = ? WHERE user_id = ?`

	_, err := db.Exec(query, newUsername, userID)
	if err != nil {
		log.Println("Error changing user username in database:", err)
		return err
	}

	return nil
}
