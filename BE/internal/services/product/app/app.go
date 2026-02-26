package app

import (
	"context"

	"golf-store/be/internal/platform/config"
	"golf-store/be/internal/platform/httpserver"
	producthttp "golf-store/be/internal/services/product/http"
	productmsg "golf-store/be/internal/services/product/messaging"
)

type App struct {
	cfg      config.ServiceConfig
	consumer *productmsg.Consumer
}

func New(cfg config.ServiceConfig) *App {
	return &App{
		cfg:      cfg,
		consumer: productmsg.NewConsumer(cfg.ServiceName, cfg.KafkaBrokers, cfg.RabbitMQURL),
	}
}

func (a *App) Run(ctx context.Context) error {
	go a.consumer.Start(ctx)

	router := producthttp.New(a.cfg.ServiceName)
	return httpserver.Run(ctx, httpserver.ServerConfig{
		Name:    a.cfg.ServiceName,
		Port:    a.cfg.HTTPPort,
		Handler: router,
	})
}
