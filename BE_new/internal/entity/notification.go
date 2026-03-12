package entity

import "time"

type Notification struct {
	ID        string     `gorm:"column:id;type:varchar(36);primaryKey"`
	UserID    string     `gorm:"column:user_id;type:varchar(36);not null;index:uk_notification_dedupe,unique"`
	Title     string     `gorm:"column:title;type:varchar(255);not null"`
	Content   string     `gorm:"column:content;type:text;not null"`
	Status    string     `gorm:"column:status;type:varchar(16);not null;default:SENT"`
	IsRead    bool       `gorm:"column:is_read;not null;default:false"`
	SentAt    *time.Time `gorm:"column:sent_at" json:"sent_at,omitempty"`
	CreatedAt time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;not null" json:"updated_at"`
	User      *User      `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}

func (Notification) TableName() string {
	return "notifications"
}
