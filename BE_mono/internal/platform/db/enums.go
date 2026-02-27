package db

const (
	RoleUser  = "USER"
	RoleAdmin = "ADMIN"
)

const (
	UserStatusActive   = "ACTIVE"
	UserStatusLocked   = "LOCKED"
	UserStatusDisabled = "DISABLED"
)

const (
	ProductStatusActive       = "ACTIVE"
	ProductStatusInactive     = "INACTIVE"
	ProductStatusDiscontinued = "DISCONTINUED"
)

const (
	OrderStatusNew            = "NEW"
	OrderStatusPendingPayment = "PENDING_PAYMENT"
	OrderStatusPaid           = "PAID"
	OrderStatusProcessing     = "PROCESSING"
	OrderStatusShipping       = "SHIPPING"
	OrderStatusCompleted      = "COMPLETED"
	OrderStatusCancelled      = "CANCELLED"
)

const (
	PaymentStatusUnpaid   = "UNPAID"
	PaymentStatusPaid     = "PAID"
	PaymentStatusFailed   = "FAILED"
	PaymentStatusRefunded = "REFUNDED"
)

const (
	PaymentTxnStatePending   = "PENDING"
	PaymentTxnStateSuccess   = "SUCCESS"
	PaymentTxnStateFailed    = "FAILED"
	PaymentTxnStateCancelled = "CANCELLED"
)

const (
	PaymentTxnTypePayment = "PAYMENT"
	PaymentTxnTypeRefund  = "REFUND"
)

const (
	NotificationStatusPending = "PENDING"
	NotificationStatusSent    = "SENT"
	NotificationStatusFailed  = "FAILED"
)
