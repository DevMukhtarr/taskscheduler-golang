package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"taskscheduler/database"
	"taskscheduler/middlewares"
	"taskscheduler/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var signInRequest models.UserSignInRequest

	err := json.NewDecoder(r.Body).Decode(&signInRequest)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	collection := database.GetCollection("users")

	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	var oldUser models.User
	err = collection.FindOne(ctx, bson.M{"username": signInRequest.Username}).Decode(&oldUser)

	if err == nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		err := bcrypt.CompareHashAndPassword([]byte(oldUser.Password), []byte(signInRequest.Password))

		if err != nil {
			http.Error(w, "Username or password do not match", http.StatusUnauthorized)
			return
		}

		token, err := middlewares.CreateToken(oldUser.ID.Hex())

		cookie := http.Cookie{
			Name:     "x-access-token",
			Value:    token,
			Path:     "/",
			MaxAge:   3600 * 24 * 7,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}

		http.SetCookie(w, &cookie)

		if err != nil {
			http.Error(w, "Failed to create token", http.StatusInternalServerError)
			return
		}

		response := models.UserResponse{
			Username: oldUser.Username,
			Token:    token,
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Invalid username or password", http.StatusInternalServerError)
		return
	}
}
