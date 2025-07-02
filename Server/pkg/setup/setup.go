package setup

import (
	"log"
	"os"

	"github.com/ArteShow/Assistant/Server/pkg/database"
)

func SetUpDatabase() {
	log_file, err := os.OpenFile("Server/log/database.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer log_file.Close()
	log.SetOutput(log_file)

	log.Println("Starting database setup...")
	db, err := database.OpenDataBase()
	if err != nil {
		log.Fatal(err)
	}
	if err := database.SetupDatabase(db); err != nil {
		log.Fatal("DB setup failed:", err)
	}

}
