package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golf-store/be/internal/platform/config"
	gatewayapp "golf-store/be/internal/services/gateway/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load("api-gateway", "8080")
	app := gatewayapp.New(cfg)

	if err := app.Run(ctx); err != nil {
		log.Fatalf("[api-gateway] fatal error: %v", err)
	}
}
