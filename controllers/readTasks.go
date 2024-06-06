package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"taskscheduler/database"
	"taskscheduler/middlewares"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ReadTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	collection := database.GetCollection("tasks")

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	user_id := r.Context().Value(middlewares.UserIDKey).(string)

	objectID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	cursor, err := collection.Find(ctx, bson.M{"user_id": objectID})

	if err != nil {
		http.Error(w, "task does not exist", http.StatusNotFound)
		return
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		http.Error(w, "result not decoded successfully", http.StatusExpectationFailed)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(results)
}
