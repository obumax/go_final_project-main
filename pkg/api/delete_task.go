package api

import (
	"net/http"

	"go1f/pkg/db"
)

// DeleteTaskHandler обрабатывает запросы на удаление задачи
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	err := db.DeleteTask(id)
	if err != nil {
		writeError(w, http.StatusOK, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{})
}
