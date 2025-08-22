package application

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ArteShow/Assistant/Server/pkg/configloader"
)

func AddTask(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to add task")

	res, err := http.Post("http://localhost:8083/internal/task/add", "application/json", r.Body)
	if err != nil {
		log.Println("Error making POST request to internal server:", err)
		http.Error(w, "Failed to forward request", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	w.WriteHeader(res.StatusCode)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	logFile, err := os.OpenFile("Server/log/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "Failed to open log file", http.StatusInternalServerError)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Println("Received request to delete task")

	res, err := http.Post("http://localhost:8083/internal/task/delete", "application/json", r.Body)
	if err != nil {
		log.Println("Error making DELETE request to internal server:", err)
		http.Error(w, "Failed to forward delete request", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	w.WriteHeader(res.StatusCode)
}

func GetTasksList(w http.ResponseWriter, r *http.Request) {
	logFile, err := os.OpenFile("Server/log/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "Failed to open log file", http.StatusInternalServerError)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.Println("Getting tasks for user")

	res, err := http.Get("http://localhost:8083/internal/task/getTasksList")
	if err != nil {
		log.Println("Failed to get Task List")
		http.Error(w, "Failed to get Task List", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.StatusCode)
	w.Write(body)
}

func GetTaskById(w http.ResponseWriter, r *http.Request) {
	logFile, err := os.OpenFile("Server/log/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Println("Getting task for u, hi-hi")

	resp, err := http.Post("http://localhost:8083/internal/getTaskByID", "application/json", r.Body)
	if err != nil {
		http.Error(w, "Error while getting your task", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error while reading Task from Internal Server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func GetAllUsersTasks(w http.ResponseWriter, r *http.Request) {
	logFile, err := os.OpenFile("Server/log/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	resp, err := http.Post("http://localhost:8083/internal/getUsersTaskList", "application/json", r.Body)
	if err != nil {
		http.Error(w, "Error while getting your task", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error while reading Tasks from Internal Server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func EditTasksStatus(w http.ResponseWriter, r *http.Request) {
	logFile, err := os.OpenFile("Server/log/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	resp, err := http.Post("http://localhost:8083/internal/editTasksStatus", "application/json", r.Body)
	if err != nil {
		http.Error(w, "Error while getting your task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
}

func AddMoney(w http.ResponseWriter, r *http.Request) {
	logFile, err := os.OpenFile("Server/log/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	resp, err := http.Post("http://localhost:8083/internal/money/addMoney", "application/json", r.Body)
	if err != nil {
		http.Error(w, "Error while setting new amount of money", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
}

func SetGoalForMoneyDatabase(w http.ResponseWriter, r *http.Request) {
	logFile, err := os.OpenFile("Server/log/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	resp, err := http.Post("http://localhost:8083/internal/money/setGoal", "application/json", r.Body)
	if err != nil {
		http.Error(w, "Error while setting new goal", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
}

func GetMoneyDatabaseStats(w http.ResponseWriter, r *http.Request) {
	logFile, err := os.OpenFile("Server/log/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "Failed to open log file", http.StatusInternalServerError)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	res, err := http.Get("http://localhost:8083/internal/money/getStats")
	if err != nil {
		log.Println("Failed to get Task List")
		http.Error(w, "Failed to get Task List", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	log.Println(body, "That was the body as you can see")
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.StatusCode)
	w.Write(body)
}

func RegistereNewUser(w http.ResponseWriter, r *http.Request) {
	logFile, err := os.OpenFile("Server/log/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	resp, err := http.Post("http://localhost:8083/internal/register/newUser", "application/json", r.Body)
	if err != nil {
		http.Error(w, "Error while creating new user", http.StatusInternalServerError)
		return
	}
	log.Println(resp.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
}

func LoginNewUser(w http.ResponseWriter, r *http.Request) {
	logFile, err := os.OpenFile("Server/log/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	resp, err := http.Post("http://localhost:8083/internal/login", "application/json", r.Body)
	if err != nil {
		http.Error(w, "Error while getting your task", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error while reading new token from internal", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func StartApplicationServer() error {
	logFile, err := os.OpenFile("Server/log/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening log file:", err)
		return err
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Println("Starting application server...")

	port, err := configloader.GetApplicationPort()
	if err != nil {
		log.Println("Error loading config:", err)
		return err
	}

	portStr := ":" + strconv.Itoa(port)
	//Tasks
	http.HandleFunc("/task/add", AddTask)
	http.HandleFunc("/task/delete", DeleteTask)
	http.HandleFunc("/task/getTasksList", GetTasksList)
	http.HandleFunc("/task/getTaskByID", GetTaskById)
	http.HandleFunc("/task/getUsersTaskList", GetAllUsersTasks)
	http.HandleFunc("/task/editTasksStatus", EditTasksStatus)

	//Money
	http.HandleFunc("/money/setGoal", SetGoalForMoneyDatabase)
	http.HandleFunc("/money/addMoney", AddMoney)
	http.HandleFunc("/money/getStats", GetMoneyDatabaseStats)

	//Registration
	http.HandleFunc("/registration", RegistereNewUser)
	http.HandleFunc("/login", LoginNewUser)

	return http.ListenAndServe(portStr, nil)
}
