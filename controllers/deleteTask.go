package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"taskscheduler/database"
	"taskscheduler/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var deleteTask models.DeleteTaskRequest

	err := json.NewDecoder(r.Body).Decode(&deleteTask)

	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
	}

	taskID, err := primitive.ObjectIDFromHex(deleteTask.TaskID)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	collection := database.GetCollection("task")

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	result := collection.FindOneAndDelete(ctx, bson.M{"_id": taskID})
	if result.Err() != nil {
		fmt.Println(err)
		http.Error(w, "not able to delete task", http.StatusExpectationFailed)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"message": "Task deleted successfully"})
}
