package application

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ArteShow/Assistant/Server/pkg/authorization"
	"github.com/ArteShow/Assistant/Server/pkg/configloader"
	"github.com/ArteShow/Assistant/Server/pkg/database"
)

// =============== Middleware & Helpers ===============

type contextKey string

var userIDKey = contextKey("userID")

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		db, err := database.OpenDataBase()
		if err != nil {
			http.Error(w, "Failed to open database", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		claims, err := authorization.ValidateJWT(tokenStr, db)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(r *http.Request) int64 {
	if val, ok := r.Context().Value(userIDKey).(int64); ok {
		return val
	}
	return 0
}

// =============== Handlers ===============

func AddTask(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to add task")

	userID := GetUserIDFromContext(r)

	// Read body into map to inject user_id
	var bodyMap map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&bodyMap); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	bodyMap["user_id"] = userID
	bodyBytes, _ := json.Marshal(bodyMap)

	res, err := http.Post("http://localhost:8083/internal/task/add", "application/json", strings.NewReader(string(bodyBytes)))
	if err != nil {
		log.Println("Error making POST request to internal server:", err)
		http.Error(w, "Failed to forward request", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	w.WriteHeader(res.StatusCode)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to delete task")

	userID := GetUserIDFromContext(r)

	var bodyMap map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&bodyMap); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	bodyMap["user_id"] = userID
	bodyBytes, _ := json.Marshal(bodyMap)

	res, err := http.Post("http://localhost:8083/internal/task/delete", "application/json", strings.NewReader(string(bodyBytes)))
	if err != nil {
		log.Println("Error making DELETE request to internal server:", err)
		http.Error(w, "Failed to forward delete request", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	w.WriteHeader(res.StatusCode)
}

func GetTasksList(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting tasks for user")

	userID := GetUserIDFromContext(r)

	bodyMap := map[string]interface{}{"user_id": userID}
	bodyBytes, _ := json.Marshal(bodyMap)

	res, err := http.Post("http://localhost:8083/internal/task/getTasksList", "application/json", strings.NewReader(string(bodyBytes)))
	if err != nil {
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
	log.Println("Getting task for user")

	userID := GetUserIDFromContext(r)

	var bodyMap map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&bodyMap); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	bodyMap["user_id"] = userID
	bodyBytes, _ := json.Marshal(bodyMap)

	resp, err := http.Post("http://localhost:8083/internal/getTaskByID", "application/json", strings.NewReader(string(bodyBytes)))
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
	log.Println("Getting all tasks for user")

	userID := GetUserIDFromContext(r)

	bodyMap := map[string]interface{}{"user_id": userID}
	bodyBytes, _ := json.Marshal(bodyMap)

	resp, err := http.Post("http://localhost:8083/internal/getUsersTaskList", "application/json", strings.NewReader(string(bodyBytes)))
	if err != nil {
		http.Error(w, "Error while getting tasks", http.StatusInternalServerError)
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
	log.Println("Editing task status")

	userID := GetUserIDFromContext(r)

	var bodyMap map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&bodyMap); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	bodyMap["user_id"] = userID
	bodyBytes, _ := json.Marshal(bodyMap)

	resp, err := http.Post("http://localhost:8083/internal/editTasksStatus", "application/json", strings.NewReader(string(bodyBytes)))
	if err != nil {
		http.Error(w, "Error while editing task status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
}

func AddMoney(w http.ResponseWriter, r *http.Request) {
	log.Println("Adding money")

	userID := GetUserIDFromContext(r)

	var bodyMap map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&bodyMap); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	bodyMap["user_id"] = userID
	bodyBytes, _ := json.Marshal(bodyMap)

	resp, err := http.Post("http://localhost:8083/internal/money/addMoney", "application/json", strings.NewReader(string(bodyBytes)))
	if err != nil {
		http.Error(w, "Error while adding money", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
}

func SetGoalForMoneyDatabase(w http.ResponseWriter, r *http.Request) {
	log.Println("Setting money goal")

	userID := GetUserIDFromContext(r)

	var bodyMap map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&bodyMap); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	bodyMap["user_id"] = userID
	bodyBytes, _ := json.Marshal(bodyMap)

	resp, err := http.Post("http://localhost:8083/internal/money/setGoal", "application/json", strings.NewReader(string(bodyBytes)))
	if err != nil {
		http.Error(w, "Error while setting new goal", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
}

func GetMoneyDatabaseStats(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting money stats")

	userID := GetUserIDFromContext(r)

	bodyMap := map[string]interface{}{"user_id": userID}
	bodyBytes, _ := json.Marshal(bodyMap)

	res, err := http.Post("http://localhost:8083/internal/money/getStats", "application/json", strings.NewReader(string(bodyBytes)))
	if err != nil {
		http.Error(w, "Failed to get stats", http.StatusInternalServerError)
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

func RegistereNewUser(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Post("http://localhost:8083/internal/register/newUser", "application/json", r.Body)
	if err != nil {
		http.Error(w, "Error while creating new user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
}

func LoginNewUser(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Post("http://localhost:8083/internal/login", "application/json", r.Body)
	if err != nil {
		http.Error(w, "Error while logging in", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error while reading token from internal", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

// =============== Server Start ===============

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

	// Tasks
	http.Handle("/task/add", JWTMiddleware(http.HandlerFunc(AddTask)))
	http.Handle("/task/delete", JWTMiddleware(http.HandlerFunc(DeleteTask)))
	http.Handle("/task/getTasksList", JWTMiddleware(http.HandlerFunc(GetTasksList)))
	http.Handle("/task/getTaskByID", JWTMiddleware(http.HandlerFunc(GetTaskById)))
	http.Handle("/task/getUsersTaskList", JWTMiddleware(http.HandlerFunc(GetAllUsersTasks)))
	http.Handle("/task/editTasksStatus", JWTMiddleware(http.HandlerFunc(EditTasksStatus)))

	// Money
	http.Handle("/money/setGoal", JWTMiddleware(http.HandlerFunc(SetGoalForMoneyDatabase)))
	http.Handle("/money/addMoney", JWTMiddleware(http.HandlerFunc(AddMoney)))
	http.Handle("/money/getStats", JWTMiddleware(http.HandlerFunc(GetMoneyDatabaseStats)))

	// Registration & login (public)
	http.HandleFunc("/registration", RegistereNewUser)
	http.HandleFunc("/login", LoginNewUser)

	return http.ListenAndServe(portStr, nil)
}
