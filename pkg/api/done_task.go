package api

import (
	"fmt"
	"net/http"
	"time"

	"go1f/pkg/db"
	"go1f/pkg/service"
)

// doneTaskHandler обрабатывает запросы на завершение задачи (/api/task/done)
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, http.StatusOK, err.Error())
		return
	}

	// Если нет правила повторения, задача удаляется
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeError(w, http.StatusOK, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{})
		return
	}

	// Если есть правило повторения, вычисляется следующая дата
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	nextDate, err := service.NextDate(today, task.Date, task.Repeat)
	if err != nil {
		writeError(w, http.StatusOK, fmt.Sprintf("ошибка при вычислении следующей даты: %v", err))
		return
	}

	if err := db.UpdateDate(nextDate, id); err != nil {
		writeError(w, http.StatusOK, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}
