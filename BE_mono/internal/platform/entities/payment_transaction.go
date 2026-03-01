package entities

import "time"

type PaymentTransaction struct {
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
	CreatedAt        time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
	Order            *Order    `gorm:"foreignKey:OrderID;references:ID" json:"order,omitempty"`
}

func (PaymentTransaction) TableName() string {
	return "payment_transactions"
}
