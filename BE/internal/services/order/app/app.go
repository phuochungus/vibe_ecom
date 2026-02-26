package app

import (
	"context"

	"golf-store/be/internal/platform/config"
	"golf-store/be/internal/platform/httpserver"
	orderhttp "golf-store/be/internal/services/order/http"
	ordermsg "golf-store/be/internal/services/order/messaging"
)

type App struct {
	cfg      config.ServiceConfig
	consumer *ordermsg.Consumer
}

func New(cfg config.ServiceConfig) *App {
	return &App{
		cfg:      cfg,
		consumer: ordermsg.NewConsumer(cfg.ServiceName, cfg.KafkaBrokers, cfg.RabbitMQURL),
	}
}

func (a *App) Run(ctx context.Context) error {
	go a.consumer.Start(ctx)

	router := orderhttp.New(a.cfg.ServiceName)
	return httpserver.Run(ctx, httpserver.ServerConfig{
		Name:    a.cfg.ServiceName,
		Port:    a.cfg.HTTPPort,
		Handler: router,
	})
}
