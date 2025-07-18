package main

import (
	"log"

	"go1f/pkg/config"
	"go1f/pkg/db"
	"go1f/pkg/server"
)

func main() {

	// Загрузка настроек
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("ошибка загрузки настроек: %v", err)
	}

	// Инициализация базы данных
	if err := db.Init(cfg.DBFile); err != nil {
		log.Fatalf("ошибка инициализации базы данных: %v", err)
	}
	defer db.Close()

	// Запуск сервера
	if err := server.Run(cfg.Port, cfg.Password); err != nil {
		log.Fatalf("ошибка запуска сервера: %v", err)
	}
}
