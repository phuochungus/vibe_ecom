package main

import (
	db "BE_new/internal/db"
	"fmt"

	"github.com/gin-gonic/gin"

	_ "github.com/joho/godotenv/autoload"
)

func main() {

	err := db.InitializeDatabase()
	if err != nil {
		panic(fmt.Errorf("failed to initialize database: %w", err))
	}
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.Run(":8080")
}
