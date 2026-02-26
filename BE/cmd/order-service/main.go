package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golf-store/be/internal/platform/config"
	orderapp "golf-store/be/internal/services/order/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load("order-service", "8083")
	app := orderapp.New(cfg)

	if err := app.Run(ctx); err != nil {
		log.Fatalf("[order-service] fatal error: %v", err)
	}
}
