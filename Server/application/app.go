package application

import (
	"errors"
	"log"
	"net/http"
	"os"
)

func CreateLogger(logFile string) *log.Logger {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(file, "", log.LstdFlags|log.Lshortfile)
	return logger
}

func NewHandler(w http.ResponseWriter, r *http.Request) {
	logger := CreateLogger("./log/server.log")
	logger.Println("New task added")
}

func StartServer(){
	logger := CreateLogger("./log/server.log")
	logger.Println("Starting server...")
	Config, err := LoadConfig("./configs/config/json")
	if err != nil {
		logger.Fatal("Error loading config:", err)
		errors.New("Failed to load config")
		panic(err)
	}
	http.HandleFunc("/task/add", NewHandler)
	http.ListenAndServe(":"+string(Config.Port), nil)
}

