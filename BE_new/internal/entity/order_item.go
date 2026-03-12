package entity

import "time"

type OrderItem struct {
	ID        string    `gorm:"column:id" json:"id"`
	OrderID   string    `gorm:"column:order_id" json:"order_id"`
	ProductID string    `gorm:"column:product_id" json:"product_id"`
	SKU       string    `gorm:"column:sku" json:"sku"`
	Name      string    `gorm:"column:name" json:"name"`
	UnitPrice int64     `gorm:"column:unit_price" json:"unit_price"`
	Quantity  int       `gorm:"column:quantity" json:"quantity"`
	LineTotal int64     `gorm:"column:line_total" json:"line_total"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Order     *Order    `json:"order,omitempty"`
	Product   *Product  `json:"product,omitempty"`
}

func (OrderItem) TableName() string {
	return "order_items"
}
