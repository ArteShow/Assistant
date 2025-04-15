package database

import (
	"database/sql"
	"log"
	"os"
)

func ConnectToDatabase(dbPath string) (*sql.DB, error) {
	log_file, err := os.OpenFile("Server/log/database.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	log.SetOutput(log_file)
	if err != nil {
		log.Println("Error opening log file:", err)
		return nil, err
	}
	
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error connecting to database:", err)
		return nil, err
	}
	
	err = db.Ping()
	if err != nil {
		log.Println("Error pinging database:", err)
		return nil, err
	}
	
	log.Println("Connected to database successfully")
	return db, nil
}

func CloseDatabase(db *sql.DB) {
	log_file, err := os.OpenFile("Server/log/database.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	log.SetOutput(log_file)
	if err != nil {
		log.Println("Error opening log file:", err)
	}

	err2 := db.Close()
	if err2 != nil {
		log.Println("Error closing database:", err)
	} else {
		log.Println("Database closed successfully")
	}
}

func AddValueToDatabase(tablename string, value any, db *sql.DB) error{
	log_file, err := os.OpenFile("Server/log/database.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	log.SetOutput(log_file)
	if err != nil {
		log.Println("Error opening log file:", err)
		return err
	}
	
	query := "INSERT INTO " + tablename + " VALUES (?)"
	_, err = db.Exec(query, value)
	if err != nil {
		log.Println("Error adding value to database:", err)
		return err
	}
	
	log.Println("Value added to database successfully")
	return nil
}

func DeletValueFromDataBase(tablename string, value any, db *sql.DB) error{
	log_file, err := os.OpenFile("Server/log/database.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	log.SetOutput(log_file)
	if err != nil {
		log.Println("Error opening log file:", err)
		return err
	}
	
	query := "DELETE FROM " + tablename + " WHERE value = ?"
	_, err = db.Exec(query, value)
	if err != nil {
		log.Println("Error deleting value from database:", err)
		return err
	}
	
	log.Println("Value deleted from database successfully")
	return nil
}