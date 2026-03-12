package routes

import (
	"BE_new/internal/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(_ *gorm.DB) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/ping", handlers.Ping)

	return router
}
