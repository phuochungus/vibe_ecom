package server

import (
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"

	authhttp "golf-store/be-mono/internal/modules/auth/http"
	notificationhttp "golf-store/be-mono/internal/modules/notification/http"
	orderhttp "golf-store/be-mono/internal/modules/order/http"
	paymenthttp "golf-store/be-mono/internal/modules/payment/http"
	producthttp "golf-store/be-mono/internal/modules/product/http"
	reportinghttp "golf-store/be-mono/internal/modules/reporting/http"
	userhttp "golf-store/be-mono/internal/modules/user/http"
	"golf-store/be-mono/internal/platform/entities"
	"golf-store/be-mono/internal/platform/observability"
	"golf-store/be-mono/internal/shared/middleware"
)

const swaggerUIHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  <style>
    html, body { margin: 0; padding: 0; }
    #swagger-ui { max-width: 1200px; margin: 0 auto; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function () {
      SwaggerUIBundle({
        url: "/openapi.yaml",
        dom_id: "#swagger-ui",
        deepLinking: true
      });
    };
  </script>
</body>
</html>`

type RouteModules struct {
	Auth          *authhttp.Handler
	Products      *producthttp.Handler
	Orders        *orderhttp.Handler
	Payments      *paymenthttp.Handler
	Notifications *notificationhttp.Handler
	Reporting     *reportinghttp.Handler
	Users         *userhttp.Handler
	AuthResolver  middleware.TokenResolver
}

func NewRouter(health *HealthHandler, modules RouteModules, allowedOrigins []string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), observability.CorrelationID(), middleware.RequestID())
	r.Use(corsMiddleware(allowedOrigins))

	r.Static("/assets", "public")
	registerDocsRoutes(r)

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
			modules.Users.Register(protected)
		}

		admin := protected.Group("/admin")
		admin.Use(middleware.RequireRole(entities.RoleAdmin))
		{
			modules.Products.RegisterAdmin(admin)
			modules.Orders.RegisterAdmin(admin)
			modules.Reporting.RegisterAdmin(admin)
		}
	}

	return r
}

func corsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	normalized := make([]string, 0, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			normalized = append(normalized, trimmed)
		}
	}

	allowAnyOrigin := slices.Contains(normalized, "*")

	return func(c *gin.Context) {
		origin := strings.TrimSpace(c.GetHeader("Origin"))
		if origin != "" {
			if allowAnyOrigin {
				c.Header("Access-Control-Allow-Origin", "*")
			} else if slices.Contains(normalized, origin) {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Vary", "Origin")
			}
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Requested-With, X-Request-Id")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func registerDocsRoutes(r *gin.Engine) {
	r.GET("/docs", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerUIHTML))
	})
	r.GET("/docs/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerUIHTML))
	})

	r.GET("/openapi.yaml", func(c *gin.Context) {
		specPath, ok := resolveOpenAPISpecPath()
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "openapi spec file not found",
				"hint":  "set OPENAPI_SPEC_PATH or place file at docs/FE/openapi.yaml",
			})
			return
		}
		c.File(specPath)
	})
}

func resolveOpenAPISpecPath() (string, bool) {
	candidates := []string{
		os.Getenv("OPENAPI_SPEC_PATH"),
		"docs/FE/openapi.yaml",
		"../docs/FE/openapi.yaml",
	}

	for _, candidate := range candidates {
		p := strings.TrimSpace(candidate)
		if p == "" {
			continue
		}
		absPath, err := filepath.Abs(p)
		if err != nil {
			continue
		}
		info, err := os.Stat(absPath)
		if err != nil || info.IsDir() {
			continue
		}
		return absPath, true
	}
	return "", false
}
