package api

import (
	"net/http"
	"time"

	"go1f/pkg/db"
)

const defaultLimit = 50

type tasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

// getTasksHandler - обработчик GET-запросов для /api/tasks
func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
		return
	}

	search := r.URL.Query().Get("search")
	dateParam := ""
	if search != "" {
		if t, err := time.Parse("02.01.2006", search); err == nil {
			dateParam = t.Format("20060102")
			search = ""
		}
	}

	tasks, err := db.Tasks(defaultLimit, search, dateParam)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Чтобы JSON не дал "tasks":null
	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}

	writeJSON(w, http.StatusOK, tasksResponse{Tasks: tasks})
}
