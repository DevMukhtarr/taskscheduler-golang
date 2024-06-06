package app

import (
	"taskscheduler/routes"
)

func App() {
	// create new task
	routes.CreateNewTask()
	routes.NewUser()
}
