package entities

import "time"

type User struct {
	ID                  string         `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	Email               string         `gorm:"column:email;type:varchar(255);uniqueIndex;not null" json:"email"`
	Phone               string         `gorm:"column:phone;type:varchar(20);uniqueIndex;not null" json:"phone"`
	Password            string         `gorm:"column:password;type:varchar(255);not null" json:"-"`
	FullName            string         `gorm:"column:full_name;type:varchar(150);not null" json:"full_name"`
	Role                string         `gorm:"column:role;type:enum('USER','ADMIN');not null;default:USER" json:"role"`
	Status              string         `gorm:"column:status;type:enum('ACTIVE','LOCKED','DISABLED');not null;default:ACTIVE" json:"status"`
	FailedLoginAttempts int            `gorm:"column:failed_login_attempts;not null;default:0" json:"failed_login_attempts"`
	LockedUntil         *time.Time     `gorm:"column:locked_until" json:"locked_until,omitempty"`
	LastLoginAt         *time.Time     `gorm:"column:last_login_at" json:"last_login_at,omitempty"`
	CreatedAt           time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
	AuthTokens          []AuthToken    `gorm:"foreignKey:UserID;references:ID" json:"-"`
	Addresses           []UserAddress  `gorm:"foreignKey:UserID;references:ID" json:"addresses,omitempty"`
	Orders              []Order        `gorm:"foreignKey:UserID;references:ID" json:"orders,omitempty"`
	Notifications       []Notification `gorm:"foreignKey:UserID;references:ID" json:"notifications,omitempty"`
}

func (User) TableName() string {
	return "users"
}
