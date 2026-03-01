package entities

import (
	"time"

	"gorm.io/datatypes"
)

type Array[T any] []T

type Product struct {
	ID          string         `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	SKU         string         `gorm:"column:sku;type:varchar(64);uniqueIndex;not null" json:"sku"`
	Name        string         `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Description string         `gorm:"column:description;type:text" json:"description"`
	Price       int64          `gorm:"column:price;not null" json:"price"`
	Stock       int            `gorm:"column:stock;not null" json:"stock"`
	Status      string         `gorm:"column:status;type:enum('ACTIVE','INACTIVE','DISCONTINUED');not null;default:ACTIVE" json:"status"`
	ImageURL    string         `gorm:"column:image_url;type:varchar(500)" json:"image_url"`
	ImageURLs   datatypes.JSON `gorm:"column:image_urls;type:json" json:"image_urls"`
	DeletedAt   *time.Time     `gorm:"column:deleted_at"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
	OrderItems  []OrderItem    `gorm:"foreignKey:ProductID;references:ID" json:"order_items,omitempty"`
	Orders      []Order        `gorm:"many2many:order_items;foreignKey:ID;joinForeignKey:ProductID;references:ID;joinReferences:OrderID" json:"orders,omitempty"`
}

func (Product) TableName() string {
	return "products"
}
