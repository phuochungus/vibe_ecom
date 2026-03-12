package entity

import "time"

type Order struct {
	ID             string      `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	OrderCode      string      `gorm:"column:order_code;type:varchar(32);uniqueIndex;not null" json:"order_code"`
	UserID         string      `gorm:"column:user_id;type:varchar(36);not null;index:uk_orders_user_idempotency,unique" json:"user_id"`
	OrderStatus    string      `gorm:"column:order_status;type:varchar(32);not null" json:"order_status"`
	SubtotalAmount int64       `gorm:"column:subtotal_amount;not null" json:"subtotal_amount"`
	DiscountAmount int64       `gorm:"column:discount_amount;not null" json:"discount_amount"`
	ShippingFee    int64       `gorm:"column:shipping_fee;not null" json:"shipping_fee"`
	TotalAmount    int64       `gorm:"column:total_amount;not null" json:"total_amount"`
	PaymentDueAt   *time.Time  `gorm:"column:payment_due_at"`
	CustomerNote   *string     `gorm:"column:customer_note;type:varchar(500)"`
	CancelReason   *string     `gorm:"column:cancel_reason;type:varchar(255)" json:"cancel_reason,omitempty"`
	FullAddress    string      `gorm:"column:full_address;type:varchar(500);not null" json:"full_address"`
	PlacedAt       time.Time   `gorm:"column:placed_at;not null" json:"placed_at"`
	CreatedAt      time.Time   `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt      time.Time   `gorm:"column:updated_at;not null" json:"updated_at"`
	User           *User       `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Items          []OrderItem `gorm:"foreignKey:OrderID;references:ID" json:"items,omitempty"`
	Products       []Product   `gorm:"many2many:order_items;foreignKey:ID;joinForeignKey:OrderID;references:ID;joinReferences:ProductID" json:"products,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}
