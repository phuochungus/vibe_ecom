package entities

import "time"

type OrderTrackingEvent struct {
	ID          string    `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	OrderID     string    `gorm:"column:order_id;type:varchar(36);not null;index" json:"order_id"`
	FromStatus  *string   `gorm:"column:from_status;type:varchar(32)" json:"from_status,omitempty"`
	ToStatus    string    `gorm:"column:to_status;type:varchar(32);not null" json:"to_status"`
	SourceType  string    `gorm:"column:source_type;type:varchar(32);not null" json:"source_type"`
	Description *string   `gorm:"column:description;type:varchar(500)" json:"description,omitempty"`
	OccurredAt  time.Time `gorm:"column:occurred_at;not null" json:"occurred_at"`
	Order       *Order    `gorm:"foreignKey:OrderID;references:ID" json:"order,omitempty"`
}

func (OrderTrackingEvent) TableName() string {
	return "order_tracking_events"
}
