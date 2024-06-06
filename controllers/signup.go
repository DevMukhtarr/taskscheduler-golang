package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"taskscheduler/database"
	"taskscheduler/middlewares"
	"taskscheduler/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Message: message})
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var signUpRequest models.UserSignUpRequest

	err := json.NewDecoder(r.Body).Decode(&signUpRequest)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if signUpRequest.Password != signUpRequest.Confirm_Password {
		http.Error(w, "Password does not match", http.StatusBadRequest)
		return
	}
	user_id := primitive.NewObjectID()
	encrypted_password, err := bcrypt.GenerateFromPassword([]byte(signUpRequest.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Fatal(err)
	}
	newUser := models.User{
		ID:       user_id,
		Username: signUpRequest.Username,
		Password: string(encrypted_password),
	}

	collection := database.GetCollection("users")

	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	// var oldUser bson.Raw
	_, err = collection.FindOne(ctx, bson.M{"username": newUser.Username}).Raw()

	if err == nil {
		sendErrorResponse(w, http.StatusConflict, "User already exists")
		return
	}

	_, err = collection.InsertOne(ctx, newUser)

	if err != nil {
		http.Error(w, "Failed to create new user", http.StatusInternalServerError)
		return
	}

	token, err := middlewares.CreateToken(user_id.Hex())

	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}
	response := models.UserResponse{
		Username: newUser.Username,
		Token:    token,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
