package api

import (
	"encoding/json"
	"net/http"

	"go1f/pkg/db"
)

// updateTaskHandler обрабатывает запросы на обновление задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeError(w, http.StatusOK, "ошибка декодирования JSON")
		return
	}
	if err := validateTask(&t); err != nil {
		writeError(w, http.StatusOK, err.Error())
		return
	}
	if err := db.UpdateTask(&t); err != nil {
		writeError(w, http.StatusOK, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{})
}
