package models

import "time"

type Notification struct {
	ID        string     `gorm:"column:id" json:"id"`
	UserID    string     `gorm:"column:user_id" json:"user_id"`
	Title     string     `gorm:"column:title" json:"title"`
	Content   string     `gorm:"column:content" json:"content"`
	Status    string     `gorm:"column:status" json:"status"`
	IsRead    bool       `gorm:"column:is_read" json:"is_read"`
	SentAt    *time.Time `gorm:"column:sent_at" json:"sent_at,omitempty"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (Notification) TableName() string {
	return "notifications"
}
