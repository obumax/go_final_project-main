package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go1f/pkg/db"
)

// updateTaskHandler обрабатывает запросы на обновление задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("ошибка декодирования JSON: %v", err))
		return
	}
	if err := validateTask(&t); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := db.UpdateTask(&t); err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("не удалось обновить задачу: %v", err))
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{})
}
