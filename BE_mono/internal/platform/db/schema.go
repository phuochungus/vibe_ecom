package db

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	entities "golf-store/be-mono/internal/platform/entities"
)

func InitSchema(gdb *gorm.DB) error {
	if gdb == nil {
		return fmt.Errorf("gorm db is required")
	}

	return gdb.AutoMigrate(
		&entities.User{},
		&entities.AuthToken{},
		&entities.UserAddress{},
		&entities.Product{},
		&entities.Order{},
		&entities.OrderItem{},
		&entities.OrderStatusHistory{},
		&entities.OrderTrackingEvent{},
		&entities.ShipmentTrackingEvent{},
		&entities.PaymentTransaction{},
		&entities.Notification{},
		&entities.AuditLog{},
	)
}

func SeedDemoData(gdb *gorm.DB) error {
	if gdb == nil {
		return fmt.Errorf("gorm db is required")
	}

	now := time.Now().UTC()
	return gdb.Transaction(func(tx *gorm.DB) error {
		if err := ensureSeedUser(tx, seedUserInput{
			Email:    "admin@golf.local",
			Phone:    "0900000001",
			Password: "admin123",
			FullName: "System Admin",
			Role:     "ADMIN",
			Status:   "ACTIVE",
			Now:      now,
		}); err != nil {
			return err
		}

		if err := ensureSeedUser(tx, seedUserInput{
			Email:    "user@golf.local",
			Phone:    "0900000002",
			Password: "user123",
			FullName: "Demo User",
			Role:     "USER",
			Status:   "ACTIVE",
			Now:      now,
		}); err != nil {
			return err
		}

		var productCount int64
		if err := tx.Model(&entities.Product{}).Count(&productCount).Error; err != nil {
			return err
		}
		if productCount == 0 {
			products := []entities.Product{
				{
					ID:          uuid.NewString(),
					SKU:         "GLF-DRIVER-001",
					Name:        "Driver Pro X",
					Description: "460cc driver for distance and forgiveness",
					Price:       129900,
					Stock:       25,
					Status:      "ACTIVE",
					ImageURL:    "https://images.local/driver-pro-x.jpg",
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:          uuid.NewString(),
					SKU:         "GLF-IRON-SET-002",
					Name:        "Iron Set Tour 6pcs",
					Description: "Forged iron set for mid-to-low handicaps",
					Price:       219900,
					Stock:       12,
					Status:      "ACTIVE",
					ImageURL:    "https://images.local/iron-tour-6pcs.jpg",
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:          uuid.NewString(),
					SKU:         "GLF-PUTTER-003",
					Name:        "Putter Classic Blade",
					Description: "Face-balanced blade putter",
					Price:       89900,
					Stock:       30,
					Status:      "ACTIVE",
					ImageURL:    "https://images.local/putter-classic-blade.jpg",
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			}
			if err := tx.Create(&products).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

type seedUserInput struct {
	Email    string
	Phone    string
	Password string
	FullName string
	Role     string
	Status   string
	Now      time.Time
}

func ensureSeedUser(tx *gorm.DB, input seedUserInput) error {
	var user entities.User
	err := tx.Where("email = ?", strings.TrimSpace(input.Email)).Take(&user).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		hashed, hashErr := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			return hashErr
		}
		return tx.Create(&entities.User{
			ID:                  uuid.NewString(),
			Email:               strings.TrimSpace(input.Email),
			Phone:               strings.TrimSpace(input.Phone),
			Password:            string(hashed),
			FullName:            strings.TrimSpace(input.FullName),
			Role:                input.Role,
			Status:              input.Status,
			FailedLoginAttempts: 0,
			CreatedAt:           input.Now,
			UpdatedAt:           input.Now,
		}).Error
	}

	updates := map[string]any{
		"full_name":  strings.TrimSpace(input.FullName),
		"role":       input.Role,
		"status":     input.Status,
		"updated_at": input.Now,
	}

	if !looksLikeBcrypt(user.Password) {
		hashed, hashErr := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			return hashErr
		}
		updates["password"] = string(hashed)
	}

	return tx.Model(&entities.User{}).Where("id = ?", user.ID).Updates(updates).Error
}

func looksLikeBcrypt(hash string) bool {
	return strings.HasPrefix(strings.TrimSpace(hash), "$2")
}
