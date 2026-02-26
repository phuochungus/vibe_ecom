package http

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"
)

func New(serviceName string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(nethttp.StatusOK, gin.H{
			"service": serviceName,
			"status":  "ok",
		})
	})

	r.GET("/readyz", func(c *gin.Context) {
		c.JSON(nethttp.StatusOK, gin.H{
			"service": serviceName,
			"ready":   true,
		})
	})

	return r
}
