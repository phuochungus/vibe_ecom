package model

import "time"

type UserRole string

type UserStatus string

type ProductStatus string

type OrderStatus string

type PaymentStatus string

type PaymentTxnState string

type PaymentTxnType string

type NotificationStatus string

const (
	RoleUser  UserRole = "USER"
	RoleAdmin UserRole = "ADMIN"
)

const (
	UserStatusActive   UserStatus = "ACTIVE"
	UserStatusLocked   UserStatus = "LOCKED"
	UserStatusDisabled UserStatus = "DISABLED"
)

const (
	ProductStatusActive       ProductStatus = "ACTIVE"
	ProductStatusInactive     ProductStatus = "INACTIVE"
	ProductStatusDiscontinued ProductStatus = "DISCONTINUED"
)

const (
	OrderStatusNew            OrderStatus = "NEW"
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusProcessing     OrderStatus = "PROCESSING"
	OrderStatusShipping       OrderStatus = "SHIPPING"
	OrderStatusCompleted      OrderStatus = "COMPLETED"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

const (
	PaymentStatusUnpaid   PaymentStatus = "UNPAID"
	PaymentStatusPaid     PaymentStatus = "PAID"
	PaymentStatusFailed   PaymentStatus = "FAILED"
	PaymentStatusRefunded PaymentStatus = "REFUNDED"
)

const (
	PaymentTxnStatePending   PaymentTxnState = "PENDING"
	PaymentTxnStateSuccess   PaymentTxnState = "SUCCESS"
	PaymentTxnStateFailed    PaymentTxnState = "FAILED"
	PaymentTxnStateCancelled PaymentTxnState = "CANCELLED"
)

const (
	PaymentTxnTypePayment PaymentTxnType = "PAYMENT"
	PaymentTxnTypeRefund  PaymentTxnType = "REFUND"
)

const (
	NotificationStatusPending NotificationStatus = "PENDING"
	NotificationStatusSent    NotificationStatus = "SENT"
	NotificationStatusFailed  NotificationStatus = "FAILED"
)

type User struct {
	ID                  string
	Email               string
	Phone               string
	Password            string
	FullName            string
	Role                UserRole
	Status              UserStatus
	FailedLoginAttempts int
	LockedUntil         *time.Time
	LastLoginAt         *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type Product struct {
	ID          string
	SKU         string
	Name        string
	Description string
	PriceCents  int64
	Stock       int
	Status      ProductStatus
	ImageURL    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type ShippingAddress struct {
	RecipientName  string
	RecipientPhone string
	Line1          string
	Line2          string
	Ward           string
	District       string
	City           string
	Province       string
	PostalCode     string
	CountryCode    string
}

type OrderItem struct {
	ID        string
	OrderID   string
	ProductID string
	SKU       string
	Name      string
	UnitPrice int64
	Quantity  int
	LineTotal int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Order struct {
	ID             string
	OrderCode      string
	UserID         string
	OrderStatus    OrderStatus
	PaymentStatus  PaymentStatus
	CurrencyCode   string
	SubtotalAmount int64
	DiscountAmount int64
	ShippingFee    int64
	TotalAmount    int64
	PaymentDueAt   *time.Time
	CustomerNote   string
	CancelReason   string
	Shipping       ShippingAddress
	PlacedAt       time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Items          []OrderItem
}

type TrackingEvent struct {
	ID          string
	OrderID     string
	FromStatus  string
	ToStatus    string
	SourceType  string
	Description string
	OccurredAt  time.Time
}

type PaymentTransaction struct {
	ID               string
	OrderID          string
	TxnType          PaymentTxnType
	Provider         string
	ProviderTxnCode  string
	IdempotencyKey   string
	Amount           int64
	CurrencyCode     string
	Status           PaymentTxnState
	ProviderResponse string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type Notification struct {
	ID        string
	UserID    string
	Channel   string
	EventType string
	EventKey  string
	Title     string
	Content   string
	Status    NotificationStatus
	Read      bool
	SentAt    *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
