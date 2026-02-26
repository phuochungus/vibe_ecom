package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golf-store/be/internal/platform/config"
	userapp "golf-store/be/internal/services/user/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load("user-service", "8081")
	app := userapp.New(cfg)

	if err := app.Run(ctx); err != nil {
		log.Fatalf("[user-service] fatal error: %v", err)
	}
}
