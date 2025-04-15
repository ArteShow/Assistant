package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/ArteShow/Assistant/Server/pkg/configloader"
)

func OpenDatabase() (*sql.DB, error) {
	log_file, err := os.OpenFile("Server/log/database.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer log_file.Close()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	log.SetOutput(log_file)

	dbPath, err := configloader.GetDatabasePath()
	if err != nil {
		log.Println("Error getting database path:", err)
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening database:", err)
		return nil, err
	}

	return db, nil
}