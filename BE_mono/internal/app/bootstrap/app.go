package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"golf-store/be-mono/internal/app/server"
	authhttp "golf-store/be-mono/internal/modules/auth/http"
	authsvc "golf-store/be-mono/internal/modules/auth/service"
	notificationhttp "golf-store/be-mono/internal/modules/notification/http"
	notificationsvc "golf-store/be-mono/internal/modules/notification/service"
	orderhttp "golf-store/be-mono/internal/modules/order/http"
	ordersvc "golf-store/be-mono/internal/modules/order/service"
	paymenthttp "golf-store/be-mono/internal/modules/payment/http"
	paymentsvc "golf-store/be-mono/internal/modules/payment/service"
	producthttp "golf-store/be-mono/internal/modules/product/http"
	productsvc "golf-store/be-mono/internal/modules/product/service"
	reportinghttp "golf-store/be-mono/internal/modules/reporting/http"
	reportingsvc "golf-store/be-mono/internal/modules/reporting/service"
	"golf-store/be-mono/internal/platform/config"
	"golf-store/be-mono/internal/platform/db"
	"golf-store/be-mono/internal/platform/httpserver"
)

type App struct {
	cfg    config.Config
	db     *gorm.DB
	router *gin.Engine
}

func New(cfg config.Config) (*App, error) {
	database, err := db.OpenMySQL(cfg.MySQLDSN)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := db.InitSchema(database); err != nil {
		closeGorm(database)
		return nil, fmt.Errorf("init schema: %w", err)
	}
	if err := db.SeedDemoData(database); err != nil {
		closeGorm(database)
		return nil, fmt.Errorf("seed demo data: %w", err)
	}

	health := server.NewHealthHandler(cfg.ServiceName, database)
	authService := authsvc.New(database)
	productService := productsvc.New(database)
	orderService := ordersvc.New(database)
	paymentService := paymentsvc.New(database, orderService)
	notificationService := notificationsvc.New(database)
	reportingService := reportingsvc.New(database)

	router := server.NewRouter(health, server.RouteModules{
		Auth:          authhttp.New(authService),
		Products:      producthttp.New(productService),
		Orders:        orderhttp.New(orderService),
		Payments:      paymenthttp.New(paymentService),
		Notifications: notificationhttp.New(notificationService),
		Reporting:     reportinghttp.New(reportingService),
		AuthResolver:  authService,
	})

	return &App{
		cfg:    cfg,
		db:     database,
		router: router,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	defer a.close()

	return httpserver.Run(ctx, httpserver.ServerConfig{
		Name:            a.cfg.ServiceName,
		Port:            a.cfg.HTTPPort,
		Handler:         a.router,
		ShutdownTimeout: 10 * time.Second,
	})
}

func (a *App) close() {
	closeGorm(a.db)
}

func closeGorm(gdb *gorm.DB) {
	if gdb == nil {
		return
	}
	sqlDB, err := gdb.DB()
	if err != nil {
		return
	}
	_ = sqlDB.Close()
}
