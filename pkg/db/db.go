package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    date    CHAR(8)   NOT NULL DEFAULT '',
    title   VARCHAR(255) NOT NULL DEFAULT '',
    comment TEXT      NOT NULL DEFAULT '',
    repeat  VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler (date);
`

var db *sql.DB

// Init открывает файл базы данных, при отсутствии создает схему
func Init(dbFile string) error {
	var err error
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка открытия базы данных %s: %w", dbFile, err)
	}
	// создаются таблицы, если их нет
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return fmt.Errorf("ошибка создания схемы: %w", err)
	}
	return nil
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
