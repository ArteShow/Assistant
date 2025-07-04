package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "modernc.org/sqlite"

	"github.com/ArteShow/Assistant/Server/pkg/configloader"
	"github.com/ArteShow/Assistant/Server/pkg/database"

	"github.com/ArteShow/Assistant/Server/pkg/task"
)

func AddTask(w http.ResponseWriter, r *http.Request) {
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
	db, err := database.OpenDataBase()
	if err != nil {
		log.Println("Error opening database:", err)
		http.Error(w, "Failed to open database", http.StatusInternalServerError)
		return
	}
	log.Println(task.SaveTask(Task, db))
	defer db.Close()
	w.WriteHeader(http.StatusOK)
}

func StartInternalServer() error {
	log.Println("Starting internal server...")

	port, err := configloader.GetInternalPort()
	if err != nil {
		return err
	}

	portStr := ":" + strconv.Itoa(port)
	http.HandleFunc("/internal/task/add", AddTask)

	return http.ListenAndServe(portStr, nil)
}
