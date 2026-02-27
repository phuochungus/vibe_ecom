package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golf-store/be-mono/internal/app/bootstrap"
	"golf-store/be-mono/internal/platform/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load("be-mono", "8080")
	app, err := bootstrap.New(cfg)
	if err != nil {
		log.Fatalf("[be-mono] bootstrap failed: %v", err)
	}

	if err := app.Run(ctx); err != nil {
		log.Fatalf("[be-mono] fatal error: %v", err)
	}
}
