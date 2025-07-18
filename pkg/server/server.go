package server

import (
	"fmt"
	"log"
	"net/http"

	"go1f/pkg/api"
	"go1f/pkg/api/middleware"
)

// Run запускает HTTP-сервер c JWT-паролем
func Run(port int, jwtSecret string) error {
	// jwtSecret передается в api/middleware
	middleware.SetTokenSecret(jwtSecret)

	// инициализация обработчиков
	if err := api.Init(); err != nil {
		return fmt.Errorf("ошибка инициализации обработчиков: %w", err)
	}

	http.Handle("/", http.FileServer(http.Dir("web")))

	address := fmt.Sprintf(":%d", port)
	log.Printf("сервер запущен на порту: %s\n", address)
	return http.ListenAndServe(address, nil)
}
