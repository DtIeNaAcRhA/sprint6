package server

import (
	"log"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/handlers"
)

// Server структура для HTTP сервера
type Server struct {
	Logger *log.Logger
	Server *http.Server
}

// New создает новый экземпляр сервера
func New(logger *log.Logger) *Server {
	// Создаем роутер
	router := http.NewServeMux()

	// Регистрируем обработчики
	router.HandleFunc("/", handlers.HandleMain)
	router.HandleFunc("/upload", handlers.HandleUpload)

	// Создаем и конфигурируем HTTP сервер
	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		Logger: logger,
		Server: httpServer,
	}
}
