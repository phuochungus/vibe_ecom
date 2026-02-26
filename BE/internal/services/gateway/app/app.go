package app

import (
	"context"

	"golf-store/be/internal/platform/config"
	"golf-store/be/internal/platform/httpserver"
	gatewayhttp "golf-store/be/internal/services/gateway/http"
	gatewaymsg "golf-store/be/internal/services/gateway/messaging"
)

type App struct {
	cfg       config.ServiceConfig
	publisher gatewaymsg.CommandPublisher
}

func New(cfg config.ServiceConfig) *App {
	return &App{
		cfg:       cfg,
		publisher: gatewaymsg.NewBrokerPublisher(cfg.ServiceName),
	}
}

func (a *App) Run(ctx context.Context) error {
	router := gatewayhttp.New(a.cfg.ServiceName, a.publisher)
	return httpserver.Run(ctx, httpserver.ServerConfig{
		Name:    a.cfg.ServiceName,
		Port:    a.cfg.HTTPPort,
		Handler: router,
	})
}
