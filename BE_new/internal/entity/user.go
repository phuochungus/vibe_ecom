package entity

import "time"

type User struct {
	ID            string         `gorm:"column:id" json:"id"`
	Email         string         `gorm:"column:email" json:"email"`
	Phone         string         `gorm:"column:phone" json:"phone"`
	Password      string         `gorm:"column:password" json:"-"`
	FullName      string         `gorm:"column:full_name" json:"full_name"`
	Role          UserRole       `gorm:"column:role" json:"role"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at" json:"updated_at"`
	Addresses     []UserAddress  `json:"addresses,omitempty"`
	Orders        []Order        `json:"orders,omitempty"`
	Notifications []Notification `json:"notifications,omitempty"`
}

func (User) TableName() string {
	return "users"
}
