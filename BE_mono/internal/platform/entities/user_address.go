package entities

import "time"

type UserAddress struct {
	ID             string     `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	UserID         string     `gorm:"column:user_id;type:varchar(36);not null;index:idx_user_addresses__user_id;index:idx_user_addresses__user_default,priority:1" json:"user_id"`
	RecipientName  string     `gorm:"column:recipient_name;type:varchar(150);not null" json:"recipient_name"`
	RecipientPhone string     `gorm:"column:recipient_phone;type:varchar(20);not null" json:"recipient_phone"`
	Line1          string     `gorm:"column:line1;type:varchar(255);not null" json:"line1"`
	Line2          *string    `gorm:"column:line2;type:varchar(255)"`
	Ward           *string    `gorm:"column:ward;type:varchar(120)"`
	District       *string    `gorm:"column:district;type:varchar(120)"`
	City           string     `gorm:"column:city;type:varchar(120);not null" json:"city"`
	Province       *string    `gorm:"column:province;type:varchar(120)" json:"province"`
	PostalCode     *string    `gorm:"column:postal_code;type:varchar(20)"`
	CountryCode    string     `gorm:"column:country_code;type:char(2);not null;default:VN"`
	IsDefault      bool       `gorm:"column:is_default;not null;default:false;index:idx_user_addresses__user_default,priority:2"`
	DeletedAt      *time.Time `gorm:"column:deleted_at"`
	CreatedAt      time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;not null" json:"updated_at"`
	User           *User      `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}

func (UserAddress) TableName() string {
	return "user_addresses"
}
