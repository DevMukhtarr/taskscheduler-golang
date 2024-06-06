package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"taskscheduler/database"
	"taskscheduler/middlewares"
	"taskscheduler/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var taskRequest models.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&taskRequest)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized: missing user ID", http.StatusUnauthorized)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if taskRequest.Hour == "" || taskRequest.Minute == "" || taskRequest.Task == "" {
		http.Error(w, "Field can not be empty", http.StatusBadRequest)
		return
	}

	task := models.Task{
		ID:     primitive.NewObjectID(),
		UserID: objectID,
		Task:   taskRequest.Task,
		Hour:   taskRequest.Hour,
		Minute: taskRequest.Minute,
	}

	collection := database.GetCollection("tasks")
	// create a new context
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, task)

	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}
