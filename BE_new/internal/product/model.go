package product

import (
	"time"

	"gorm.io/datatypes"
)

type Product struct {
	ID          string         `gorm:"column:id" json:"id"`
	SKU         string         `gorm:"column:sku" json:"sku"`
	Name        string         `gorm:"column:name" json:"name"`
	Description string         `gorm:"column:description" json:"description"`
	Price       int64          `gorm:"column:price" json:"price"`
	Stock       int            `gorm:"column:stock" json:"stock"`
	Status      string         `gorm:"column:status" json:"status"`
	ImageURL    string         `gorm:"column:image_url" json:"image_url"`
	ImageURLs   datatypes.JSON `gorm:"column:image_urls" json:"image_urls"`
	DeletedAt   *time.Time     `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
}

func (Product) TableName() string {
	return "products"
}
