package entities

import "time"

type Order struct {
	ID                    string                  `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	OrderCode             string                  `gorm:"column:order_code;type:varchar(32);uniqueIndex;not null" json:"order_code"`
	IdempotencyKey        string                  `gorm:"column:idempotency_key;type:varchar(128);not null;index:uk_orders_user_idempotency,unique" json:"idempotency_key"`
	UserID                string                  `gorm:"column:user_id;type:varchar(36);not null;index:uk_orders_user_idempotency,unique" json:"user_id"`
	OrderStatus           string                  `gorm:"column:order_status;type:varchar(32);not null" json:"order_status"`
	PaymentStatus         string                  `gorm:"column:payment_status;type:varchar(16);not null" json:"payment_status"`
	CurrencyCode          string                  `gorm:"column:currency_code;type:char(3);not null" json:"currency_code"`
	SubtotalAmount        int64                   `gorm:"column:subtotal_amount;not null" json:"subtotal_amount"`
	DiscountAmount        int64                   `gorm:"column:discount_amount;not null" json:"discount_amount"`
	ShippingFee           int64                   `gorm:"column:shipping_fee;not null" json:"shipping_fee"`
	TotalAmount           int64                   `gorm:"column:total_amount;not null" json:"total_amount"`
	PaymentDueAt          *time.Time              `gorm:"column:payment_due_at"`
	CustomerNote          *string                 `gorm:"column:customer_note;type:varchar(500)"`
	CancelReason          *string                 `gorm:"column:cancel_reason;type:varchar(255)" json:"cancel_reason,omitempty"`
	ShippingRecipientName string                  `gorm:"column:shipping_recipient_name;type:varchar(150);not null" json:"shipping_recipient_name"`
	ShippingPhone         string                  `gorm:"column:shipping_phone;type:varchar(20);not null" json:"shipping_phone"`
	ShippingLine1         string                  `gorm:"column:shipping_line1;type:varchar(255);not null" json:"shipping_line1"`
	ShippingLine2         *string                 `gorm:"column:shipping_line2;type:varchar(255)" json:"shipping_line2,omitempty"`
	ShippingWard          *string                 `gorm:"column:shipping_ward;type:varchar(120)" json:"shipping_ward,omitempty"`
	ShippingDistrict      *string                 `gorm:"column:shipping_district;type:varchar(120)" json:"shipping_district,omitempty"`
	ShippingCity          string                  `gorm:"column:shipping_city;type:varchar(120);not null" json:"shipping_city"`
	ShippingProvince      *string                 `gorm:"column:shipping_province;type:varchar(120)" json:"shipping_province,omitempty"`
	ShippingPostalCode    *string                 `gorm:"column:shipping_postal_code;type:varchar(20)" json:"shipping_postal_code,omitempty"`
	ShippingCountryCode   string                  `gorm:"column:shipping_country_code;type:char(2);not null" json:"shipping_country_code"`
	PlacedAt              time.Time               `gorm:"column:placed_at;not null" json:"placed_at"`
	CreatedAt             time.Time               `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt             time.Time               `gorm:"column:updated_at;not null" json:"updated_at"`
	User                  *User                   `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Items                 []OrderItem             `gorm:"foreignKey:OrderID;references:ID" json:"items,omitempty"`
	Products              []Product               `gorm:"many2many:order_items;foreignKey:ID;joinForeignKey:OrderID;references:ID;joinReferences:ProductID" json:"products,omitempty"`
	StatusHistories       []OrderStatusHistory    `gorm:"foreignKey:OrderID;references:ID" json:"status_histories,omitempty"`
	ShipmentEvents        []ShipmentTrackingEvent `gorm:"foreignKey:OrderID;references:ID" json:"shipment_events,omitempty"`
	OrderTrackingEvents   []OrderTrackingEvent    `gorm:"foreignKey:OrderID;references:ID" json:"order_tracking_events,omitempty"`
	PaymentTransactions   []PaymentTransaction    `gorm:"foreignKey:OrderID;references:ID" json:"payment_transactions,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}
