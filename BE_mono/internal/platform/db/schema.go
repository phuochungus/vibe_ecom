package db

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	entities "golf-store/be-mono/internal/platform/entities"
)

const seedUserPassword = "123456"

//go:embed taylormade_products.json
var taylormadeProductsJSON []byte

var legacyDemoProductSKUs = []string{
	"GLF-DRIVER-001",
	"GLF-IRON-SET-002",
	"GLF-PUTTER-003",
}

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

func SeedDemoData(gdb *gorm.DB, publicBaseURL string) error {
	if gdb == nil {
		return fmt.Errorf("gorm db is required")
	}

	now := time.Now().UTC()
	return gdb.Transaction(func(tx *gorm.DB) error {
		if err := ensureSeedUser(tx, seedUserInput{
			Email:    "admin@golf.local",
			Phone:    "0900000001",
			Password: seedUserPassword,
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
			Password: seedUserPassword,
			FullName: "Demo User",
			Role:     "USER",
			Status:   "ACTIVE",
			Now:      now,
		}); err != nil {
			return err
		}

		if err := retireLegacyDemoProducts(tx, now); err != nil {
			return err
		}

		products, err := loadTaylorMadeSeedProducts(now)
		if err != nil {
			return err
		}

		for _, product := range products {
			if err := ensureSeedProduct(tx, product); err != nil {
				return err
			}
		}

		return nil
	})
}

type seedProductRecord struct {
	SKU         string `json:"sku"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Stock       int    `json:"stock"`
	Status      string `json:"status"`
	ImageURL    string `json:"image_url"`
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

func seedImageURL(publicBaseURL string, filename string) string {
	baseURL := strings.TrimRight(strings.TrimSpace(publicBaseURL), "/")
	if baseURL == "" {
		baseURL = "http://127.0.0.1:8080"
	}
	return baseURL + "/assets/seed-images/" + strings.TrimSpace(filename)
}

func loadTaylorMadeSeedProducts(now time.Time) ([]entities.Product, error) {
	var records []seedProductRecord
	if err := json.Unmarshal(taylormadeProductsJSON, &records); err != nil {
		return nil, fmt.Errorf("decode taylormade products: %w", err)
	}

	products := make([]entities.Product, 0, len(records))
	for _, record := range records {
		product := entities.Product{
			ID:          uuid.NewString(),
			SKU:         strings.TrimSpace(record.SKU),
			Name:        strings.TrimSpace(record.Name),
			Description: strings.TrimSpace(record.Description),
			Price:       record.Price,
			Stock:       record.Stock,
			Status:      normalizeSeedProductStatus(record.Status),
			ImageURL:    strings.TrimSpace(record.ImageURL),
			ImageURLs:   datatypes.JSON([]byte("[]")),
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		products = append(products, product)
	}

	return products, nil
}

func normalizeSeedProductStatus(value string) string {
	if strings.EqualFold(strings.TrimSpace(value), entities.ProductStatusInactive) {
		return entities.ProductStatusInactive
	}
	return entities.ProductStatusActive
}

func retireLegacyDemoProducts(tx *gorm.DB, now time.Time) error {
	return tx.Model(&entities.Product{}).
		Where("sku IN ?", legacyDemoProductSKUs).
		Updates(map[string]any{
			"status":     entities.ProductStatusInactive,
			"deleted_at": now,
			"updated_at": now,
		}).Error
}

func ensureSeedProduct(tx *gorm.DB, input entities.Product) error {
	var product entities.Product
	err := tx.Where("sku = ?", strings.TrimSpace(input.SKU)).Take(&product).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return tx.Create(&input).Error
	}

	updates := map[string]any{
		"name":        strings.TrimSpace(input.Name),
		"description": strings.TrimSpace(input.Description),
		"price":       input.Price,
		"stock":       input.Stock,
		"status":      input.Status,
		"image_url":   strings.TrimSpace(input.ImageURL),
		"image_urls":  input.ImageURLs,
		"deleted_at":  nil,
		"updated_at":  input.UpdatedAt,
	}

	return tx.Model(&entities.Product{}).Where("id = ?", product.ID).Updates(updates).Error
}
