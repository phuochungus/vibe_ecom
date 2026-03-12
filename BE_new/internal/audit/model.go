package audit

import (
	"time"

	"gorm.io/datatypes"
)

type AuditLog struct {
	ID         string            `gorm:"column:id" json:"id"`
	EntityType string            `gorm:"column:entity_type" json:"entity_type"`
	EntityID   string            `gorm:"column:entity_id" json:"entity_id"`
	Action     string            `gorm:"column:action" json:"action"`
	ActorID    *string           `gorm:"column:actor_id" json:"actor_id,omitempty"`
	BeforeData datatypes.JSONMap `gorm:"column:before_data" json:"before_data,omitempty"`
	AfterData  datatypes.JSONMap `gorm:"column:after_data" json:"after_data,omitempty"`
	Metadata   datatypes.JSONMap `gorm:"column:metadata" json:"metadata,omitempty"`
	IPAddress  *string           `gorm:"column:ip_address" json:"ip_address,omitempty"`
	UserAgent  *string           `gorm:"column:user_agent" json:"user_agent,omitempty"`
	CreatedAt  time.Time         `gorm:"column:created_at" json:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
