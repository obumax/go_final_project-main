package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

// Значения по умолчанию
const (
	defaultPort         = 7540
	defaultNameOfDBFile = "scheduler.db"
)

// Config содержит внешние параметры
type Config struct {
	DBFile   string // Путь к файлу БД (TODO_DBFILE)
	Port     int    // Порт HTTP-сервера (TODO_PORT)
	Password string // Секрет для JWT (TODO_PASSWORD)
}

// Load читает .env и окружение, заполняет Config
// Если какие-то переменные не заданы, подставляются значения по умолчанию
func Load() (*Config, error) {
	_ = godotenv.Load()

	dbFile := getenv("TODO_DBFILE", filepath.Join(".", defaultNameOfDBFile))

	portStr := getenv("TODO_PORT", strconv.Itoa(defaultPort))
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("неверный TODO_PORT=%q: %w", portStr, err)
	}

	pass := os.Getenv("TODO_PASSWORD")

	return &Config{
		DBFile:   dbFile,
		Port:     port,
		Password: pass,
	}, nil
}

// getenv возвращает os.Getenv(key), если не пусто, иначе по умолчанию
func getenv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
