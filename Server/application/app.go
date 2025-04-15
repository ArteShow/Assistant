package application

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ArteShow/Assistant/Server/pkg/configloader"
)



func AddTask(w http.ResponseWriter, r *http.Request) {
	log_file, err := os.OpenFile("Server/log/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer log_file.Close()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to open log file", http.StatusInternalServerError)
		return
	}
	log.SetOutput(log_file)
	
	log.Println("Received request to add task")
	_, err2 := http.Post("http://localhost:8082/task/add", "application/json", r.Body)
	if err2 != nil {
		log.Println("Error making POST request:", err)
		http.Error(w, "Failed to make POST request", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(r.Response.StatusCode)
}

func StartServer(){
	log.Println("Starting server...")
	// Load the configuration file
	port, err := configloader.GetApplicationPort()
	if err != nil {
		log.Fatal("Error loading config:", err)
		errors.New("Failed to load config")
		panic(err)
	}

	port2 := ":" + strconv.Itoa(port) 
	http.HandleFunc("/task/add", AddTask)
	log.Fatal(http.ListenAndServe(port2, nil))
}

