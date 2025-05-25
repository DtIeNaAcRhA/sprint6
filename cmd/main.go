package main

import (
	"log"
	"os"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server" // замените на ваш реальный путь к пакету
)

func main() {

	logger := log.New(os.Stdout, "APP: ", log.LstdFlags|log.Lshortfile)

	srv := server.New(logger)

	logger.Println("Starting server on :8080...")
	if err := srv.Server.ListenAndServe(); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
