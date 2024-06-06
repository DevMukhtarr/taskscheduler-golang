package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"taskscheduler/database"
	"taskscheduler/middlewares"
	"taskscheduler/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var TaskUpdateRequest models.TaskUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&TaskUpdateRequest)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	collection := database.GetCollection("tasks")

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

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

	taskObjectID, err := primitive.ObjectIDFromHex(TaskUpdateRequest.TaskID)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	filter := bson.M{
		"_id":     taskObjectID,
		"user_id": objectID,
	}

	update := bson.M{
		"$set": bson.M{},
	}

	if TaskUpdateRequest.Task != "" {
		update["$set"].(bson.M)["task"] = TaskUpdateRequest.Task
	}
	if TaskUpdateRequest.Hour != "" {
		update["$set"].(bson.M)["hour"] = TaskUpdateRequest.Hour
	}
	if TaskUpdateRequest.Minute != "" {
		update["$set"].(bson.M)["minute"] = TaskUpdateRequest.Minute
	}

	if len(update["$set"].(bson.M)) == 0 {
		http.Error(w, "No valid fields to update", http.StatusBadRequest)
		return
	}

	var updatedDocument models.Task

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err = collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedDocument)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "No document was found with the provided filters", http.StatusNotFound)
		} else {
			fmt.Println(err)
		}
	} else {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Task updated successfully",
			"task":    updatedDocument.Task,
			"hour":    updatedDocument.Hour,
			"minute":  updatedDocument.Minute,
		})
	}
}
