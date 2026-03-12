package entity

import "time"

type UserAddress struct {
	ID             string     `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	UserID         string     `gorm:"column:user_id;type:varchar(36);not null;index:idx_user_addresses__user_id;index:idx_user_addresses__user_default,priority:1" json:"user_id"`
	RecipientName  string     `gorm:"column:recipient_name;type:varchar(150);not null" json:"recipient_name"`
	RecipientPhone string     `gorm:"column:recipient_phone;type:varchar(20);not null" json:"recipient_phone"`
	FullAddress    string     `gorm:"column:full_address;type:varchar(500);not null" json:"full_address"`
	IsDefault      bool       `gorm:"column:is_default;not null;default:false;index:idx_user_addresses__user_default,priority:2"`
	DeletedAt      *time.Time `gorm:"column:deleted_at"`
	CreatedAt      time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;not null" json:"updated_at"`
	User           *User      `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}

func (UserAddress) TableName() string {
	return "user_addresses"
}
