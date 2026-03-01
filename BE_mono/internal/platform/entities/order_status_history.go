package entities

import "time"

type OrderStatusHistory struct {
	ID            string    `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	OrderID       string    `gorm:"column:order_id;type:varchar(36);not null;index:idx_osh__order_time,priority:1" json:"order_id"`
	FromStatus    *string   `gorm:"column:from_status;type:varchar(32)" json:"from_status,omitempty"`
	ToStatus      string    `gorm:"column:to_status;type:varchar(32);not null;index:idx_osh__to_status_time,priority:1" json:"to_status"`
	ChangedByType string    `gorm:"column:changed_by_type;type:varchar(32);not null" json:"changed_by_type"`
	ChangedByID   *string   `gorm:"column:changed_by_id;type:varchar(36)" json:"changed_by_id,omitempty"`
	ChangeReason  *string   `gorm:"column:change_reason;type:varchar(255)" json:"change_reason,omitempty"`
	OccurredAt    time.Time `gorm:"column:occurred_at;not null;index:idx_osh__order_time,priority:2;index:idx_osh__to_status_time,priority:2" json:"occurred_at"`
	CreatedAt     time.Time `gorm:"column:created_at;not null" json:"created_at"`
	Order         *Order    `gorm:"foreignKey:OrderID;references:ID" json:"order,omitempty"`
}

func (OrderStatusHistory) TableName() string {
	return "order_status_history"
}
