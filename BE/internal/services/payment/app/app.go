package app

import (
	"context"

	"golf-store/be/internal/platform/config"
	"golf-store/be/internal/platform/httpserver"
	paymenthttp "golf-store/be/internal/services/payment/http"
	paymentmsg "golf-store/be/internal/services/payment/messaging"
)

type App struct {
	cfg      config.ServiceConfig
	consumer *paymentmsg.Consumer
}

func New(cfg config.ServiceConfig) *App {
	return &App{
		cfg:      cfg,
		consumer: paymentmsg.NewConsumer(cfg.ServiceName, cfg.KafkaBrokers, cfg.RabbitMQURL),
	}
}

func (a *App) Run(ctx context.Context) error {
	go a.consumer.Start(ctx)

	router := paymenthttp.New(a.cfg.ServiceName)
	return httpserver.Run(ctx, httpserver.ServerConfig{
		Name:    a.cfg.ServiceName,
		Port:    a.cfg.HTTPPort,
		Handler: router,
	})
}
