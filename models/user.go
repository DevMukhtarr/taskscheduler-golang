package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represents a user in the MongoDB database
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
}

type UserSignUpRequest struct {
	Username         string `json:"username"`
	Password         string `json:"password"`
	Confirm_Password string `json:"confirm_password"`
}

type UserResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type UserSignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
