package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golf-store/be/internal/platform/config"
	productapp "golf-store/be/internal/services/product/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load("product-service", "8082")
	app := productapp.New(cfg)

	if err := app.Run(ctx); err != nil {
		log.Fatalf("[product-service] fatal error: %v", err)
	}
}
