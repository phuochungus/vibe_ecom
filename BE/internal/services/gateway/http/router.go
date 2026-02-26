package http

import (
	"github.com/gin-gonic/gin"

	"golf-store/be/internal/platform/observability"
	gatewaymsg "golf-store/be/internal/services/gateway/messaging"
)

func New(serviceName string, publisher gatewaymsg.CommandPublisher) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), observability.CorrelationID())

	h := NewHandler(serviceName, publisher)

	r.GET("/healthz", h.Health)
	r.GET("/readyz", h.Ready)

	api := r.Group("/api/v1")
	{
		api.POST("/auth/login", h.Login)
		api.GET("/products", h.ListProducts)
		api.POST("/orders", h.CreateOrder)
		api.GET("/orders/:orderCode/tracking", h.GetOrderTracking)
		api.GET("/notifications", h.ListNotifications)
		api.POST("/webhooks/payment/:provider", h.ReceivePaymentWebhook)

		admin := api.Group("/admin")
		admin.GET("/orders", h.AdminListOrders)
		admin.GET("/revenue/summary", h.AdminRevenueSummary)
	}

	return r
}
