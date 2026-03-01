package entities

import "time"

type OrderItem struct {
	ID        string    `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	OrderID   string    `gorm:"column:order_id;type:varchar(36);not null;index:uk_order_item,unique" json:"order_id"`
	ProductID string    `gorm:"column:product_id;type:varchar(36);not null;index:uk_order_item,unique" json:"product_id"`
	SKU       string    `gorm:"column:sku;type:varchar(64);not null" json:"sku"`
	Name      string    `gorm:"column:name;type:varchar(255);not null" json:"name"`
	UnitPrice int64     `gorm:"column:unit_price;not null" json:"unit_price"`
	Quantity  int       `gorm:"column:quantity;not null" json:"quantity"`
	LineTotal int64     `gorm:"column:line_total;not null" json:"line_total"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
	Order     *Order    `gorm:"foreignKey:OrderID;references:ID" json:"order,omitempty"`
	Product   *Product  `gorm:"foreignKey:ProductID;references:ID" json:"product,omitempty"`
}

func (OrderItem) TableName() string {
	return "order_items"
}
