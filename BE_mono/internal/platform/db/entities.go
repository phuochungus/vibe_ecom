package db

import (
	"time"
)

type UserEntity struct {
	ID                  string     `gorm:"column:id;type:varchar(36);primaryKey"`
	Email               string     `gorm:"column:email;type:varchar(255);uniqueIndex;not null"`
	Phone               string     `gorm:"column:phone;type:varchar(20);uniqueIndex;not null"`
	Password            string     `gorm:"column:password;type:varchar(255);not null"`
	FullName            string     `gorm:"column:full_name;type:varchar(150);not null"`
	Role                string     `gorm:"column:role;type:enum('USER','ADMIN');not null;default:USER"`
	Status              string     `gorm:"column:status;type:enum('ACTIVE','LOCKED','DISABLED');not null;default:ACTIVE"`
	FailedLoginAttempts int        `gorm:"column:failed_login_attempts;not null;default:0"`
	LockedUntil         *time.Time `gorm:"column:locked_until"`
	LastLoginAt         *time.Time `gorm:"column:last_login_at"`
	CreatedAt           time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt           time.Time  `gorm:"column:updated_at;not null"`
}

func (UserEntity) TableName() string {
	return "users"
}

type AuthTokenEntity struct {
	AccessToken  string    `gorm:"column:access_token;type:varchar(512);primaryKey"`
	RefreshToken string    `gorm:"column:refresh_token;type:varchar(512);uniqueIndex;not null"`
	UserID       string    `gorm:"column:user_id;type:varchar(36);index;not null"`
	ExpiresAt    time.Time `gorm:"column:expires_at;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;not null"`
}

func (AuthTokenEntity) TableName() string {
	return "auth_tokens"
}

type ProductEntity struct {
	ID          string     `gorm:"column:id;type:varchar(36);primaryKey"`
	SKU         string     `gorm:"column:sku;type:varchar(64);uniqueIndex;not null"`
	Name        string     `gorm:"column:name;type:varchar(255);not null"`
	Description string     `gorm:"column:description;type:text"`
	Price       int64      `gorm:"column:price_cents;not null"`
	Stock       int        `gorm:"column:stock;not null"`
	Status      string     `gorm:"column:status;type:enum('ACTIVE','INACTIVE','DISCONTINUED');not null;default:ACTIVE"`
	ImageURL    string     `gorm:"column:image_url;type:varchar(500)"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;not null"`
}

func (ProductEntity) TableName() string {
	return "products"
}

type OrderEntity struct {
	ID                    string            `gorm:"column:id;type:varchar(36);primaryKey"`
	OrderCode             string            `gorm:"column:order_code;type:varchar(32);uniqueIndex;not null"`
	IdempotencyKey        string            `gorm:"column:idempotency_key;type:varchar(128);not null;index:uk_orders_user_idempotency,unique"`
	UserID                string            `gorm:"column:user_id;type:varchar(36);not null;index:uk_orders_user_idempotency,unique"`
	OrderStatus           string            `gorm:"column:order_status;type:enum('NEW','PENDING_PAYMENT','PAID','PROCESSING','SHIPPING','COMPLETED','CANCELLED');not null"`
	PaymentStatus         string            `gorm:"column:payment_status;type:enum('UNPAID','PAID','FAILED','REFUNDED');not null"`
	CurrencyCode          string            `gorm:"column:currency_code;type:char(3);not null"`
	SubtotalAmount        int64             `gorm:"column:subtotal_amount;not null"`
	DiscountAmount        int64             `gorm:"column:discount_amount;not null"`
	ShippingFee           int64             `gorm:"column:shipping_fee;not null"`
	TotalAmount           int64             `gorm:"column:total_amount;not null"`
	PaymentDueAt          *time.Time        `gorm:"column:payment_due_at"`
	CustomerNote          *string           `gorm:"column:customer_note;type:varchar(500)"`
	CancelReason          *string           `gorm:"column:cancel_reason;type:varchar(255)"`
	ShippingRecipientName string            `gorm:"column:shipping_recipient_name;type:varchar(150);not null"`
	ShippingPhone         string            `gorm:"column:shipping_phone;type:varchar(20);not null"`
	ShippingLine1         string            `gorm:"column:shipping_line1;type:varchar(255);not null"`
	ShippingLine2         *string           `gorm:"column:shipping_line2;type:varchar(255)"`
	ShippingWard          *string           `gorm:"column:shipping_ward;type:varchar(120)"`
	ShippingDistrict      *string           `gorm:"column:shipping_district;type:varchar(120)"`
	ShippingCity          string            `gorm:"column:shipping_city;type:varchar(120);not null"`
	ShippingProvince      *string           `gorm:"column:shipping_province;type:varchar(120)"`
	ShippingPostalCode    *string           `gorm:"column:shipping_postal_code;type:varchar(20)"`
	ShippingCountryCode   string            `gorm:"column:shipping_country_code;type:char(2);not null"`
	PlacedAt              time.Time         `gorm:"column:placed_at;not null"`
	CreatedAt             time.Time         `gorm:"column:created_at;not null"`
	UpdatedAt             time.Time         `gorm:"column:updated_at;not null"`
	Items                 []OrderItemEntity `gorm:"foreignKey:OrderID;references:ID"`
}

func (OrderEntity) TableName() string {
	return "orders"
}

type OrderItemEntity struct {
	ID        string       `gorm:"column:id;type:varchar(36);primaryKey"`
	OrderID   string       `gorm:"column:order_id;type:varchar(36);not null;index:uk_order_item,unique"`
	ProductID string       `gorm:"column:product_id;type:varchar(36);not null;index:uk_order_item,unique"`
	SKU       string       `gorm:"column:sku;type:varchar(64);not null"`
	Name      string       `gorm:"column:name;type:varchar(255);not null"`
	UnitPrice int64        `gorm:"column:unit_price;not null"`
	Quantity  int          `gorm:"column:quantity;not null"`
	LineTotal int64        `gorm:"column:line_total;not null"`
	CreatedAt time.Time    `gorm:"column:created_at;not null"`
	UpdatedAt time.Time    `gorm:"column:updated_at;not null"`
	Order     *OrderEntity `gorm:"foreignKey:OrderID;references:ID"`
}

func (OrderItemEntity) TableName() string {
	return "order_items"
}

type OrderTrackingEventEntity struct {
	ID          string    `gorm:"column:id;type:varchar(36);primaryKey"`
	OrderID     string    `gorm:"column:order_id;type:varchar(36);not null;index"`
	FromStatus  *string   `gorm:"column:from_status;type:varchar(32)"`
	ToStatus    string    `gorm:"column:to_status;type:varchar(32);not null"`
	SourceType  string    `gorm:"column:source_type;type:varchar(32);not null"`
	Description *string   `gorm:"column:description;type:varchar(500)"`
	OccurredAt  time.Time `gorm:"column:occurred_at;not null"`
}

func (OrderTrackingEventEntity) TableName() string {
	return "order_tracking_events"
}

type PaymentTransactionEntity struct {
	ID               string    `gorm:"column:id;type:varchar(36);primaryKey"`
	OrderID          string    `gorm:"column:order_id;type:varchar(36);not null;index:uk_payment_order_idempotency,unique"`
	TxnType          string    `gorm:"column:txn_type;type:enum('PAYMENT','REFUND');not null;default:PAYMENT"`
	Provider         string    `gorm:"column:provider;type:varchar(50);not null"`
	ProviderTxnCode  *string   `gorm:"column:provider_txn_code;type:varchar(128);uniqueIndex"`
	IdempotencyKey   *string   `gorm:"column:idempotency_key;type:varchar(128);index:uk_payment_order_idempotency,unique"`
	Amount           int64     `gorm:"column:amount;not null"`
	CurrencyCode     string    `gorm:"column:currency_code;type:char(3);not null"`
	Status           string    `gorm:"column:status;type:enum('PENDING','SUCCESS','FAILED','CANCELLED');not null"`
	ProviderResponse *string   `gorm:"column:provider_response;type:varchar(500)"`
	CreatedAt        time.Time `gorm:"column:created_at;not null"`
	UpdatedAt        time.Time `gorm:"column:updated_at;not null"`
}

func (PaymentTransactionEntity) TableName() string {
	return "payment_transactions"
}

type NotificationEntity struct {
	ID        string     `gorm:"column:id;type:varchar(36);primaryKey"`
	UserID    string     `gorm:"column:user_id;type:varchar(36);not null;index:uk_notification_dedupe,unique"`
	Channel   string     `gorm:"column:channel;type:enum('IN_APP','EMAIL','SMS');not null;default:IN_APP"`
	EventType string     `gorm:"column:event_type;type:varchar(64);not null;index:uk_notification_dedupe,unique"`
	EventKey  string     `gorm:"column:event_key;type:varchar(128);not null;index:uk_notification_dedupe,unique"`
	Title     string     `gorm:"column:title;type:varchar(255);not null"`
	Content   string     `gorm:"column:content;type:text;not null"`
	Status    string     `gorm:"column:status;type:enum('PENDING','SENT','FAILED');not null;default:SENT"`
	IsRead    bool       `gorm:"column:is_read;not null;default:false"`
	SentAt    *time.Time `gorm:"column:sent_at"`
	CreatedAt time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt time.Time  `gorm:"column:updated_at;not null"`
}

func (NotificationEntity) TableName() string {
	return "notifications"
}
