package main

import (
	"context"
	"log"

	"github.com/maxviazov/dolina-flower-order-backend/internal/app"
)

func main() {
	ctx := context.Background()

	application := app.New()
	if err := application.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	if err := application.Run(ctx); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}
