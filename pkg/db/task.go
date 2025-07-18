package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

const DateFormat = "20060102"

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет новую задачу в базу банных
func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)

	if err != nil {
		return 0, fmt.Errorf("ошибка добавления задачи: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения ID задачи: %w", err)
	}
	return id, nil
}

// GetTask получает задачу по ID из базы данных
func GetTask(id string) (*Task, error) {
	if id == "" {
		return nil, errors.New("ID задачи не указан")
	}

	t := &Task{}
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("задача не найдена")
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка получения задачи: %w", err)
	}
	return t, nil
}

// UpdateTask обновляет задачу в базе данных по ID
func UpdateTask(t *Task) error {
	if t.ID == "" {
		return errors.New("ID задачи не указан")
	}

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat, t.ID)

	if err != nil {
		return fmt.Errorf("ошибка обновления задачи: %w", err)
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества обновленных строк: %w", err)
	}
	if cnt == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

// DeleteTask удаляет задачу из базы данных по ID
func DeleteTask(id string) error {
	if id == "" {
		return errors.New("ID задачи не указан")
	}

	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := db.Exec(query, id)

	if err != nil {
		return fmt.Errorf("ошибка удаления задачи: %w", err)
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка проверки удаления: %w", err)
	}
	if cnt == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

// UpdateDate обновляет дату задачи по ID
func UpdateDate(newDate, id string) error {
	if id == "" {
		return errors.New("ID задачи не указан")
	}

	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := db.Exec(query, newDate, id)

	if err != nil {
		return fmt.Errorf("ошибка обновления даты задачи: %w", err)
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества обновленных строк: %w", err)
	}
	if cnt == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

// Tasks выводит список задач в зависимости от параметров
func Tasks(limit int, search, dateFilter string) ([]*Task, error) {
	switch {
	case dateFilter != "":
		return TasksByDate(limit, dateFilter)
	case search != "":
		return SearchTasks(limit, search)
	default:
		return ListTasks(limit)
	}
}

// getTasks извлекает задачи из базы данных
func getTasks(query string, args ...any) ([]*Task, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	var list []*Task
	for rows.Next() {
		t := new(Task)
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строк: %w", err)
		}
		list = append(list, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка после сканирования строк: %w", err)
	}
	return list, nil
}

// ListTasks выводит ограниченный limit список задач без фильтров
func ListTasks(limit int) ([]*Task, error) {
	query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		ORDER BY date ASC
		LIMIT ?
	`
	return getTasks(query, limit)
}

// SearchTasks ищет задачу по заголовку и комментарию
func SearchTasks(limit int, keyword string) ([]*Task, error) {
	tasks, err := getTasks(`
        SELECT id, date, title, comment, repeat
        FROM scheduler
        ORDER BY date ASC
    `)
	if err != nil {
		return nil, err
	}

	// Убираем чувствительность к регистру в SQLite
	lowerKeyword := strings.ToLower(keyword)
	var result []*Task

	for _, t := range tasks {
		title := strings.ToLower(t.Title)
		comment := strings.ToLower(t.Comment)

		if strings.Contains(title, lowerKeyword) || strings.Contains(comment, lowerKeyword) {
			result = append(result, t)
		}
	}

	if len(result) > limit {
		result = result[:limit]
	}

	if result == nil {
		result = make([]*Task, 0)
	}
	return result, nil
}

// TasksByDate отбирает задачи по дате
func TasksByDate(limit int, dateFilter string) ([]*Task, error) {
	if _, err := time.Parse(DateFormat, dateFilter); err != nil {
		return nil, fmt.Errorf("неверный формат даты %q", dateFilter)
	}
	q := `
        SELECT id, date, title, comment, repeat
        FROM scheduler
        WHERE date = ?
        ORDER BY date ASC
        LIMIT ?
    `
	return getTasks(q, dateFilter, limit)
}
