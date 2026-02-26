package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golf-store/be/internal/platform/config"
	notificationapp "golf-store/be/internal/services/notification/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load("notification-service", "8085")
	app := notificationapp.New(cfg)

	if err := app.Run(ctx); err != nil {
		log.Fatalf("[notification-service] fatal error: %v", err)
	}
}
