package app

import (
	"context"

	"golf-store/be/internal/platform/config"
	"golf-store/be/internal/platform/httpserver"
	notificationhttp "golf-store/be/internal/services/notification/http"
	notificationmsg "golf-store/be/internal/services/notification/messaging"
)

type App struct {
	cfg      config.ServiceConfig
	consumer *notificationmsg.Consumer
}

func New(cfg config.ServiceConfig) *App {
	return &App{
		cfg:      cfg,
		consumer: notificationmsg.NewConsumer(cfg.ServiceName, cfg.KafkaBrokers, cfg.RabbitMQURL),
	}
}

func (a *App) Run(ctx context.Context) error {
	go a.consumer.Start(ctx)

	router := notificationhttp.New(a.cfg.ServiceName)
	return httpserver.Run(ctx, httpserver.ServerConfig{
		Name:    a.cfg.ServiceName,
		Port:    a.cfg.HTTPPort,
		Handler: router,
	})
}
