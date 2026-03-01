package dto

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
	Items    []OrderItemDTO `json:"items"`
	Payments []any          `json:"payments"`
}