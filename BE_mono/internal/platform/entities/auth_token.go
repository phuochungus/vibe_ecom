package entities

import "time"

type AuthToken struct {
	AccessToken  string    `gorm:"column:access_token;type:varchar(512);primaryKey" json:"access_token"`
	RefreshToken string    `gorm:"column:refresh_token;type:varchar(512);uniqueIndex;not null" json:"refresh_token"`
	UserID       string    `gorm:"column:user_id;type:varchar(36);index;not null" json:"user_id"`
	ExpiresAt    time.Time `gorm:"column:expires_at;not null" json:"expires_at"`
	CreatedAt    time.Time `gorm:"column:created_at;not null" json:"created_at"`
	User         *User     `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}

func (AuthToken) TableName() string {
	return "auth_tokens"
}
