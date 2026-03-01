package entities

import "time"

type ShipmentTrackingEvent struct {
	ID               string    `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	OrderID          string    `gorm:"column:order_id;type:varchar(36);not null;index:idx_ship_evt__order_time,priority:1"`
	TrackingNo       *string   `gorm:"column:tracking_no;type:varchar(64);index:idx_ship_evt__tracking_time,priority:1"`
	CarrierCode      *string   `gorm:"column:carrier_code;type:varchar(32)" json:"carrier_code,omitempty"`
	EventStatus      string    `gorm:"column:event_status;type:varchar(32);not null" json:"event_status"`
	EventDescription *string   `gorm:"column:event_description;type:varchar(500)" json:"event_description,omitempty"`
	EventTime        time.Time `gorm:"column:event_time;not null;index:idx_ship_evt__order_time,priority:2;index:idx_ship_evt__tracking_time,priority:2" json:"event_time"`
	SourceType       string    `gorm:"column:source_type;type:varchar(32);not null" json:"source_type"`
	SourceRef        *string   `gorm:"column:source_ref;type:varchar(128)" json:"source_ref,omitempty"`
	CreatedAt        time.Time `gorm:"column:created_at;not null" json:"created_at"`
	Order            *Order    `gorm:"foreignKey:OrderID;references:ID" json:"order,omitempty"`
}

func (ShipmentTrackingEvent) TableName() string {
	return "shipment_tracking_events"
}
