package application

import (
	"errors"
	"log"
	"net/http"
	"strconv"
)



func NewHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("New task added")
}

func StartServer(){
	log.Println("Starting server...")
	// Load the configuration file
	Config, err := LoadConfig("./Server/configs/config.json")
	if err != nil {
		log.Fatal("Error loading config:", err)
		errors.New("Failed to load config")
		panic(err)
	}

	port := ":" + strconv.Itoa(Config.Port) 
	log.Println(port)
	http.HandleFunc("/task/add", NewHandler)
	log.Fatal(http.ListenAndServe("8083", nil))
}

