package main

import (
	"context"
	"log"
	"os"

	"dolina-flower-order-backend/internal/app"
)

func main() {
	// Создаем контекст для приложения
	ctx := context.Background()

	// Создаем и инициализируем приложение
	app := app.New()
	if err := app.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Запускаем приложение
	if err := app.Run(ctx); err != nil {
		log.Fatalf("Application failed: %v", err)
	}

	os.Exit(0)
}
