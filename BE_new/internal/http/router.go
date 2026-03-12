package http

import (
	stdhttp "net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(_ *gorm.DB) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(stdhttp.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return router
}
