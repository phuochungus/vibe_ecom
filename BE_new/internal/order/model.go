package order

import "time"

type Order struct {
	ID             string      `gorm:"column:id" json:"id"`
	OrderCode      string      `gorm:"column:order_code" json:"order_code"`
	UserID         string      `gorm:"column:user_id" json:"user_id"`
	OrderStatus    string      `gorm:"column:order_status" json:"order_status"`
	SubtotalAmount int64       `gorm:"column:subtotal_amount" json:"subtotal_amount"`
	DiscountAmount int64       `gorm:"column:discount_amount" json:"discount_amount"`
	ShippingFee    int64       `gorm:"column:shipping_fee" json:"shipping_fee"`
	TotalAmount    int64       `gorm:"column:total_amount" json:"total_amount"`
	PaymentDueAt   *time.Time  `gorm:"column:payment_due_at" json:"payment_due_at,omitempty"`
	CustomerNote   *string     `gorm:"column:customer_note" json:"customer_note,omitempty"`
	CancelReason   *string     `gorm:"column:cancel_reason" json:"cancel_reason,omitempty"`
	FullAddress    string      `gorm:"column:full_address" json:"full_address"`
	PlacedAt       time.Time   `gorm:"column:placed_at" json:"placed_at"`
	CreatedAt      time.Time   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time   `gorm:"column:updated_at" json:"updated_at"`
	Items          []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}
