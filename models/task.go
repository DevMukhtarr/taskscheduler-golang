package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represents a user in the MongoDB database
type Task struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID primitive.ObjectID `bson:"user_id"`
	Task   string             `bson:"task"`
	Hour   string             `bson:"hour"`
	Minute string             `bson:"minute"`
}

type TaskRequest struct {
	UserID string `json:"user_id"`
	Task   string `json:"task"`
	Hour   string `json:"hour"`
	Minute string `json:"minute"`
}

type DeleteTaskRequest struct {
	TaskID string `json:"task_id"`
}

type TaskUpdateRequest struct {
	TaskID string `json:"task_id"`
	Task   string `json:"task"`
	Hour   string `json:"hour"`
	Minute string `json:"minute"`
}
