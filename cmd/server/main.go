package main

import (
	"context"
	"log"

	"github.com/maxviazov/dolina-flower-order-backend/internal/app"
)

func main() {
	ctx := context.Background()

	app := app.New()
	if err := app.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	if err := app.Run(ctx); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}
