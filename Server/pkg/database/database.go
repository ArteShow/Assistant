package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/ArteShow/Assistant/Server/pkg/configloader"
	_ "modernc.org/sqlite"
)

func SetupDatabase(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS tasks (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            description TEXT,
            status TEXT,
			money INTEGER NOT NULL,
            user_id INTEGER
        );`,
		`CREATE TABLE IF NOT EXISTS jwt_token (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			jwt_key BLOB NOT NULL
		);`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}
	log.Println("All tables created successfully.")

	return nil
}

func OpenDataBase() (*sql.DB, error) {
	logFile, err := os.OpenFile("Server/log/database.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	log.SetOutput(logFile)

	dbPath, err := configloader.GetDatabasePath()
	if err != nil {
		log.Println("Config error:", err)
		return nil, err
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Println("DB open error:", err)
		return nil, err
	}

	// Enable foreign keys!
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Println("Failed to enable foreign key constraints:", err)
		return nil, err
	}

	return db, nil
}
