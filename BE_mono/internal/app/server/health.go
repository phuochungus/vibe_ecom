package server

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"golf-store/be-mono/internal/shared/response"
)

type HealthHandler struct {
	serviceName string
	db          *sql.DB
}

func NewHealthHandler(serviceName string, db *sql.DB) *HealthHandler {
	return &HealthHandler{serviceName: serviceName, db: db}
}

func (h *HealthHandler) Health(c *gin.Context) {
	response.OK(c, gin.H{
		"service": h.serviceName,
		"status":  "ok",
	})
}

func (h *HealthHandler) Ready(c *gin.Context) {
	if h.db != nil {
		if err := h.db.PingContext(c.Request.Context()); err != nil {
			response.Error(c, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "database not ready", nil)
			return
		}
	}

	response.OK(c, gin.H{
		"service": h.serviceName,
		"ready":   true,
	})
}
