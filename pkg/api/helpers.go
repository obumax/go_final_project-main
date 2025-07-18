package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go1f/pkg/db"
	"go1f/pkg/service"
)

const DateFormat = "20060102"

// writeJSON устанавливает заголовок и записывает данные в формате JSON в ответ
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; caharset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// writeError записывает ошибку в формате JSON в ответ
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// validateTask проверяет корректность задачи перед обновлением или добавлением
func validateTask(t *db.Task) error {
	// Проверка обязательных полей
	if t.Date == "" {
		return errors.New("не указана дата задачи")
	}
	var dstart string
	if t.Date == "today" {
		dstart = time.Now().Format(DateFormat)
	} else {
		if _, err := time.Parse(DateFormat, t.Date); err != nil {
			return fmt.Errorf("неверный формат даты: %v", err)
		}
		dstart = t.Date
	}
	if strings.TrimSpace(t.Title) == "" {
		return errors.New("не указан заголовок задачи")
	}
	// Проверка правила повторения
	if t.Repeat != "" {
		if _, err := service.NextDate(time.Now(), dstart, t.Repeat); err != nil {
			return fmt.Errorf("некорректное правило повторения: %v", err)
		}
	}
	return nil
}
