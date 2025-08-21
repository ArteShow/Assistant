package internal

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "modernc.org/sqlite"

	"github.com/ArteShow/Assistant/Server/pkg/configloader"
	"github.com/ArteShow/Assistant/Server/pkg/database"

	"github.com/ArteShow/Assistant/Server/models"
	"github.com/ArteShow/Assistant/Server/pkg/task"
)

type Delete_Task struct {
	Task_ID int64 `json:"task_ID"`
	User_ID int64 `json:"user_ID"`
}

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

func GetAllUsersTasks(w http.ResponseWriter, r *http.Request) {
	log_file, err := os.OpenFile("Server/log/internal.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	defer log_file.Close()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to open log file", http.StatusInternalServerError)
		return
	}
	log.SetOutput(log_file)

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error while reading request", http.StatusInternalServerError)
		return
	}

	var IDS models.TaskIdFromRequest
	err = json.Unmarshal(body, &IDS)
	if err != nil {
		http.Error(w, "Error unmarshling the body", http.StatusInternalServerError)
		return
	}
	db, err := database.OpenDataBase()
	if err != nil {
		http.Error(w, "Failed to open database", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	log.SetOutput(log_file)

	Tasks, err := task.GetAllUsersTasks(db, IDS.User_ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to get Tasks", http.StatusInternalServerError)
		return
	}

	jsonFormatedData, err := json.Marshal(Tasks)
	if err != nil {
		http.Error(w, "Failed to formate Task into json", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonFormatedData)
}

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	log_file, err := os.OpenFile("Server/log/internal.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "Failed to open log file", http.StatusInternalServerError)
		return
	}
	defer log_file.Close()
	log.SetOutput(log_file)

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error while reading request", http.StatusInternalServerError)
		return
	}

	var IDS models.TaskIdFromRequest
	err = json.Unmarshal(body, &IDS)
	if err != nil {
		http.Error(w, "Error unmarshling the body", http.StatusInternalServerError)
		return
	}
	db, err := database.OpenDataBase()
	if err != nil {
		http.Error(w, "Failed to open database", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	log.SetOutput(log_file)

	task, err := task.GetUsersTaskByID(db, IDS.User_ID, IDS.ID)
	if err != nil {
		http.Error(w, "Error while getting Task", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	jsonFormatedData, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "Failed to formate Task into json", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonFormatedData)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	log_file, err := os.OpenFile("Server/log/internal.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "Failed to open log file", http.StatusInternalServerError)
		return
	}
	defer log_file.Close()
	log.SetOutput(log_file)

	var Delet Delete_Task
	err = json.NewDecoder(r.Body).Decode(&Delet)
	if err != nil {
		log.Println("Failed to decode the body:", err)
		http.Error(w, "Failed to decode the body", http.StatusInternalServerError)
		return
	}

	db, err := database.OpenDataBase()
	if err != nil {
		log.Println("Database error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	ok, err := task.DeleteUsersTaskByID(db, Delet.User_ID, Delet.Task_ID)
	if err != nil || !ok {
		log.Println("Failed to delete task:", err)
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetTasksList(w http.ResponseWriter, r *http.Request) {
	logFile, err := os.OpenFile("Server/log/internal.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	db, err := database.OpenDataBase()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}

	log.Println("Yo from internal!")

	list, err := task.GetTasksList(db)
	if err != nil {
		http.Error(w, "Failed to get Tasks", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(list)
	if err != nil {
		http.Error(w, "Failed to encode tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func StartInternalServer() error {
	log_file, err := os.OpenFile("Server/log/internal.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
		return err
	}
	log.SetOutput(log_file)
	log.Println("Starting internal server...")
	defer log_file.Close()
	port, err := configloader.GetInternalPort()
	if err != nil {
		return err
	}
	portStr := ":" + strconv.Itoa(port)
	http.HandleFunc("/internal/task/add", AddTask)
	http.HandleFunc("/internal/task/delete", DeleteTask)
	http.HandleFunc("/internal/task/getTasksList", GetTasksList)
	http.HandleFunc("/internal/getTaskByID", GetTaskByID)
	http.HandleFunc("/internal/getUsersTaskList", GetAllUsersTasks)

	return http.ListenAndServe(portStr, nil)
}
