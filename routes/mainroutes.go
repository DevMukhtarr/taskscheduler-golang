package routes

import (
	"net/http"
	"taskscheduler/controllers"
	"taskscheduler/middlewares"
)

func CreateNewTask() {
	http.Handle("/task/new", middlewares.CheckToken(http.HandlerFunc(controllers.CreateTask)))
	http.Handle("/task/delete", middlewares.CheckToken(http.HandlerFunc(controllers.DeleteTask)))
	http.Handle("/tasks", middlewares.CheckToken(http.HandlerFunc(controllers.ReadTasks)))
	http.Handle("/task/update", middlewares.CheckToken(http.HandlerFunc(controllers.UpdateTask)))
}
