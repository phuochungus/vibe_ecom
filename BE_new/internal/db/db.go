package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	entity "BE_new/internal/entity"
)

var DB *gorm.DB

func InitializeDatabase() error {
	POSTGRES_DSN := os.Getenv("POSTGRES_DSN")
	if POSTGRES_DSN == "" {
		return fmt.Errorf("POSTGRES_DSN environment variable is required")
	}

	dbTemp, err := gorm.Open(postgres.Open(POSTGRES_DSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = dbTemp
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(0)

	DB.AutoMigrate(&entity.User{})
	DB.AutoMigrate(&entity.UserAddress{})
	DB.AutoMigrate(&entity.Order{})
	DB.AutoMigrate(&entity.OrderItem{})
	DB.AutoMigrate(&entity.Notification{})
	DB.AutoMigrate(&entity.AuditLog{})
	DB.AutoMigrate(&entity.Product{})
	DB.AutoMigrate(&entity.Notification{})

	return nil
}
