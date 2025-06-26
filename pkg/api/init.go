package api

import (
	"net/http"

	"go1f/pkg/api/middleware"
)

// Init регистрирует все HTTP-эндпойнты
func Init() error {
	http.HandleFunc("/api/signin", signInHandler)
	http.HandleFunc("/api/nextdate", nextDateHandler)

	auth := middleware.Auth
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(getTasksHandler))
	http.HandleFunc("/api/task/done", auth(doneTaskHandler))
	return nil
}
