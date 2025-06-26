package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go1f/pkg/db"
	"go1f/pkg/service"
)

// addTaskHandler обрабатывает POST-запрос для добавления новой задачи
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// Десериализация JSON
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "ошибка декодирования JSON")
		return
	}

	// Проверка обязательного поля Title
	if task.Title == "" {
		writeError(w, http.StatusBadRequest, "не указан заголовок задачи")
		return
	}

	// Проверка и корректировка формата даты и правила повторения
	if err := checkDateFormat(&task); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Добавление задачи в базу данных
	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "ошибка добавления задачи")
		return
	}

	// Ответ с ID новой задачи
	writeJSON(w, http.StatusCreated, map[string]any{
		"error": "",
		"id":    id,
	})
}

// checkDateFormat проверяет корректность даты и правило повторения
func checkDateFormat(task *db.Task) error {
	now := time.Now()
	today := now.Format(DateFormat)

	// Если дата не указана, устанавливается текущая дата
	if task.Date == "" {
		task.Date = today
	}

	// Проверка формата даты
	if _, err := time.Parse(DateFormat, task.Date); err != nil {
		return fmt.Errorf("неверный формат даты: %v", err)
	}

	// Если нет правило повторения, устанавливается текущая дата, если она в прошлом
	if task.Repeat == "" {
		if task.Date < today {
			task.Date = today
		}
		return nil
	}

	// Если есть правило повторения, вычисляется следующая дата только если дата не сегодня
	if task.Date < today {
		next, err := service.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("ошибка вычисления следующей даты: %v", err)
		}
		task.Date = next
	}
	return nil
}
