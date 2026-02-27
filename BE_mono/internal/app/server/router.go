package server

import (
	"github.com/gin-gonic/gin"

	authhttp "golf-store/be-mono/internal/modules/auth/http"
	notificationhttp "golf-store/be-mono/internal/modules/notification/http"
	orderhttp "golf-store/be-mono/internal/modules/order/http"
	paymenthttp "golf-store/be-mono/internal/modules/payment/http"
	producthttp "golf-store/be-mono/internal/modules/product/http"
	reportinghttp "golf-store/be-mono/internal/modules/reporting/http"
	"golf-store/be-mono/internal/platform/observability"
	"golf-store/be-mono/internal/shared/middleware"
	"golf-store/be-mono/internal/shared/model"
)

type RouteModules struct {
	Auth          *authhttp.Handler
	Products      *producthttp.Handler
	Orders        *orderhttp.Handler
	Payments      *paymenthttp.Handler
	Notifications *notificationhttp.Handler
	Reporting     *reportinghttp.Handler
	AuthResolver  middleware.TokenResolver
}

func NewRouter(health *HealthHandler, modules RouteModules) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), observability.CorrelationID(), middleware.RequestID())

	r.GET("/healthz", health.Health)
	r.GET("/readyz", health.Ready)

	api := r.Group("/api/v1")
	{
		modules.Auth.RegisterPublic(api)
		modules.Products.RegisterPublic(api)
		modules.Payments.RegisterWebhook(api)

		protected := api.Group("")
		protected.Use(middleware.RequireAuth(modules.AuthResolver))
		{
			modules.Auth.RegisterProtected(protected)
			modules.Orders.RegisterUser(protected)
			modules.Payments.RegisterUser(protected)
			modules.Notifications.RegisterUser(protected)
		}

		admin := protected.Group("/admin")
		admin.Use(middleware.RequireRole(model.RoleAdmin))
		{
			modules.Products.RegisterAdmin(admin)
			modules.Orders.RegisterAdmin(admin)
			modules.Reporting.RegisterAdmin(admin)
		}
	}

	return r
}
