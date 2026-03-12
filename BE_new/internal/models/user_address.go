package models

import "time"

type UserAddress struct {
	ID             string     `gorm:"column:id" json:"id"`
	UserID         string     `gorm:"column:user_id" json:"user_id"`
	RecipientName  string     `gorm:"column:recipient_name" json:"recipient_name"`
	RecipientPhone string     `gorm:"column:recipient_phone" json:"recipient_phone"`
	FullAddress    string     `gorm:"column:full_address" json:"full_address"`
	IsDefault      bool       `gorm:"column:is_default" json:"is_default"`
	DeletedAt      *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (UserAddress) TableName() string {
	return "user_addresses"
}
