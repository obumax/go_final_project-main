package api

import (
	"net/http"

	"go1f/pkg/db"
)

// getTaskHandler обрабатывает запросы на получение задачи по ID
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, http.StatusOK, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, task)
}
