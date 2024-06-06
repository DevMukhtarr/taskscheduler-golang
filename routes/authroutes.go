package routes

import (
	"net/http"
	"taskscheduler/controllers"
)

func NewUser() {
	http.HandleFunc("/user/new", controllers.SignUp)
	http.HandleFunc("/user/signin", controllers.SignIn)
}
