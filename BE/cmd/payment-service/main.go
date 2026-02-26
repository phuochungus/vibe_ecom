package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golf-store/be/internal/platform/config"
	paymentapp "golf-store/be/internal/services/payment/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load("payment-service", "8084")
	app := paymentapp.New(cfg)

	if err := app.Run(ctx); err != nil {
		log.Fatalf("[payment-service] fatal error: %v", err)
	}
}
