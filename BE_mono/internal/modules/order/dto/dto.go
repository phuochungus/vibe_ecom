package dto

import (
	"golf-store/be-mono/internal/platform/entities"
	"time"
)

type CreateOrderItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type ShippingAddressRequest struct {
	RecipientName  string `json:"recipient_name"`
	RecipientPhone string `json:"recipient_phone"`
	Line1          string `json:"line1"`
	Line2          string `json:"line2"`
	Ward           string `json:"ward"`
	District       string `json:"district"`
	City           string `json:"city"`
	Province       string `json:"province"`
	PostalCode     string `json:"postal_code"`
	CountryCode    string `json:"country_code"`
}

type CreateOrderRequest struct {
	Items           []CreateOrderItemRequest `json:"items"`
	ShippingAddress ShippingAddressRequest   `json:"shipping_address"`
	CustomerNote    string                   `json:"customer_note"`
}

type CancelOrderRequest struct {
	Reason string `json:"reason"`
}

type AdminUpdateStatusRequest struct {
	ToStatus string `json:"to_status"`
	Reason   string `json:"reason"`
}

type OrderSummaryDTO struct {
	ID             string `json:"id"`
	OrderCode      string `json:"order_code"`
	OrderStatus    string `json:"order_status"`
	PaymentStatus  string `json:"payment_status"`
	SubtotalAmount string `json:"subtotal_amount"`
	DiscountAmount string `json:"discount_amount"`
	ShippingFee    string `json:"shipping_fee"`
	TotalAmount    string `json:"total_amount"`
	ShippingRecipientName string `json:"shipping_recipient_name,omitempty"`
	ShippingPhone         string `json:"shipping_phone,omitempty"`
	PlacedAt       string `json:"placed_at"`
	PaymentDueAt   string `json:"payment_due_at,omitempty"`
}

type OrderItemDTO struct {
	ProductID string `json:"product_id"`
	SKU       string `json:"sku"`
	Name      string `json:"name"`
	UnitPrice string `json:"unit_price"`
	Quantity  int    `json:"quantity"`
	LineTotal string `json:"line_total"`
}

type OrderDetailDTO struct {
	OrderSummaryDTO
	CustomerNote        string         `json:"customer_note,omitempty"`
	CancelReason        string         `json:"cancel_reason,omitempty"`
	ShippingLine1       string         `json:"shipping_line1"`
	ShippingLine2       string         `json:"shipping_line2,omitempty"`
	ShippingWard        string         `json:"shipping_ward,omitempty"`
	ShippingDistrict    string         `json:"shipping_district,omitempty"`
	ShippingCity        string         `json:"shipping_city"`
	ShippingProvince    string         `json:"shipping_province,omitempty"`
	ShippingPostalCode  string         `json:"shipping_postal_code,omitempty"`
	ShippingCountryCode string         `json:"shipping_country_code"`
	Items               []OrderItemDTO `json:"items"`
	Payments            []any          `json:"payments"`
}

type ShippingAddress struct {
	RecipientName  string
	RecipientPhone string
	Line1          string
	Line2          string
	Ward           string
	District       string
	City           string
	Province       string
	PostalCode     string
	CountryCode    string
}

type CreateOrderInput struct {
	UserID         string
	IdempotencyKey string
	Items          []CreateOrderItemInput
	Shipping       ShippingAddress
	CustomerNote   string
}

type CreateOrderItemInput struct {
	ProductID string
	Quantity  int
}

type ListInput struct {
	UserID   string
	Status   string
	From     *time.Time
	To       *time.Time
	Page     int
	PageSize int
	Admin    bool
}

type ListOutput struct {
	Items      []*entities.Order
	Page       int
	PageSize   int
	Total      int
	TotalPages int
}
