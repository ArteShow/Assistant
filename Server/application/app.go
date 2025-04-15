package application

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
)



func LogCreater(w http.ResponseWriter, r *http.Request) {
	logFile, err := os.OpenFile("Server/log/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "Error opening log file", http.StatusInternalServerError)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.Println("Request received:", r.Method, r.URL.Path)
}

func StartServer(){
	log.Println("Starting server...")
	// Load the configuration file
	Config, err := LoadConfig("Server/configs/config.json")
	if err != nil {
		log.Fatal("Error loading config:", err)
		errors.New("Failed to load config")
		panic(err)
	}

	port := ":" + strconv.Itoa(Config.Port) 
	http.HandleFunc("/task/add", LogCreater)
	log.Fatal(http.ListenAndServe(port, nil))
}

