package app

import (
	"context"

	"golf-store/be/internal/platform/config"
	"golf-store/be/internal/platform/httpserver"
	userhttp "golf-store/be/internal/services/user/http"
	usermsg "golf-store/be/internal/services/user/messaging"
)

type App struct {
	cfg      config.ServiceConfig
	consumer *usermsg.Consumer
}

func New(cfg config.ServiceConfig) *App {
	return &App{
		cfg:      cfg,
		consumer: usermsg.NewConsumer(cfg.ServiceName, cfg.KafkaBrokers, cfg.RabbitMQURL),
	}
}

func (a *App) Run(ctx context.Context) error {
	go a.consumer.Start(ctx)

	router := userhttp.New(a.cfg.ServiceName)
	return httpserver.Run(ctx, httpserver.ServerConfig{
		Name:    a.cfg.ServiceName,
		Port:    a.cfg.HTTPPort,
		Handler: router,
	})
}
