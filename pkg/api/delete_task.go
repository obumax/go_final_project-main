package api

import (
	"net/http"

	"go1f/pkg/db"
)

// DeleteTaskHandler обрабатывает запросы на удаление задачи
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "не передан id задачи")
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeError(w, http.StatusNotFound, "задача не найдена")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}
