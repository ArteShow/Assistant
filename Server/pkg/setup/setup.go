package setup

import (
	"log"
	"os"

	"github.com/ArteShow/Assistant/Server/pkg/database"
	money_database "github.com/ArteShow/Assistant/Server/pkg/money"
)

func SetUpDatabase() {
	log_file, err := os.OpenFile("Server/log/setup.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

	moeyDb, err := money_database.OpenDataBase()
	if err != nil {
		log.Fatal(err)
	}
	if err := money_database.SetupDatabase(moeyDb); err != nil {
		log.Fatal("DB money setup failed:", err)
	}
}
