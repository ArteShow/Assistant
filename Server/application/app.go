package application

import (
	"errors"
	"log"
	"net/http"
)



func NewHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("New task added")
}

func StartServer(){
	log.Println("Starting server...")
	// Load the configuration file
	Config, err := LoadConfig("./Server/configs/config/json")
	if err != nil {
		log.Fatal("Error loading config:", err)
		errors.New("Failed to load config")
		panic(err)
	}
	http.HandleFunc("/task/add", NewHandler)
	http.ListenAndServe(":"+string(Config.Port), nil)
}

