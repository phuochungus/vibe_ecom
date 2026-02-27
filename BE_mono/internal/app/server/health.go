package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"golf-store/be-mono/internal/shared/response"
)

type HealthHandler struct {
	serviceName string
	db          *gorm.DB
}

func NewHealthHandler(serviceName string, db *gorm.DB) *HealthHandler {
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
		sqlDB, err := h.db.DB()
		if err != nil {
			response.Error(c, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "database not ready", nil)
			return
		}
		if err := sqlDB.PingContext(c.Request.Context()); err != nil {
			response.Error(c, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "database not ready", nil)
			return
		}
	}

	response.OK(c, gin.H{
		"service": h.serviceName,
		"ready":   true,
	})
}
