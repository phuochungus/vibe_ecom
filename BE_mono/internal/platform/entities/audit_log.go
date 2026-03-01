package entities

import (
	"time"
)

type AuditLog struct {
	ID         string    `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	EntityType string    `gorm:"column:entity_type;type:varchar(64);not null;index:idx_audit__entity_time,priority:1" json:"entity_type"`
	EntityID   string    `gorm:"column:entity_id;type:varchar(36);not null;index:idx_audit__entity_time,priority:2" json:"entity_id"`
	Action     string    `gorm:"column:action;type:varchar(64);not null;index:idx_audit__action_time,priority:1" json:"action"`
	ActorType  string    `gorm:"column:actor_type;type:varchar(32);not null;index:idx_audit__actor_time,priority:1" json:"actor_type"`
	ActorID    *string   `gorm:"column:actor_id;type:varchar(36);index:idx_audit__actor_time,priority:2" json:"actor_id,omitempty"`
	BeforeData *string   `gorm:"column:before_data;type:json" json:"before_data,omitempty"`
	AfterData  *string   `gorm:"column:after_data;type:json" json:"after_data,omitempty"`
	Metadata   *string   `gorm:"column:metadata;type:json" json:"metadata,omitempty"`
	IPAddress  *string   `gorm:"column:ip_address;type:varchar(45)" json:"ip_address,omitempty"`
	UserAgent  *string   `gorm:"column:user_agent;type:varchar(255)" json:"user_agent,omitempty"`
	CreatedAt  time.Time `gorm:"column:created_at;not null;index:idx_audit__entity_time,priority:3;index:idx_audit__action_time,priority:2" json:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
