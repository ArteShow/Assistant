package internal

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"

	"github.com/ArteShow/Assistant/Server/pkg/configloader"
	"github.com/ArteShow/Assistant/Server/pkg/task"
)



func AddTask(w http.ResponseWriter, r *http.Request){
	log_file, err := os.OpenFile("Server/log/internal.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer log_file.Close()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to open log file", http.StatusInternalServerError)
		return
	}
	log.SetOutput(log_file)

	Task := &task.Task{}
	err = json.NewDecoder(r.Body).Decode(Task)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	dbpath, err := configloader.GetDatabasePath()
	if err != nil {
		log.Println("Error opening database:", err)
		http.Error(w, "Failed to open database", http.StatusInternalServerError)
		return
	}
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		log.Println("Error opening database:", err)
		http.Error(w, "Failed to open database", http.StatusInternalServerError)
		return
	}
	log.Println(task.SaveTask(Task, Task.ID, db))
	defer db.Close()
	w.WriteHeader(http.StatusOK)
}

func StartServer(){
	log.Println("Starting server internal")
	// Load the configuration file
	port, err := configloader.GetInternalPort()
	if err != nil {
		log.Fatal("Error loading config:", err)
		panic(err)
	}

	port2 := ":" + strconv.Itoa(port) 
	http.HandleFunc("/internal/task/add", AddTask)
	http.ListenAndServe(port2, nil)
}