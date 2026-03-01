package dto

type SummaryResponseDTO struct {
	From               string  `json:"from"`
	To                 string  `json:"to"`
	GrossRevenue       string  `json:"gross_revenue"`
	RefundAmount       string  `json:"refund_amount"`
	NetRevenue         string  `json:"net_revenue"`
	CompletedOrders    int     `json:"completed_orders"`
	PaymentSuccessRate string  `json:"payment_success_rate"`
}

type RevenueOrderItemDTO struct {
	ID             string `json:"id"`
	OrderCode      string `json:"order_code"`
	OrderStatus    string `json:"order_status"`
	PaymentStatus  string `json:"payment_status"`
	SubtotalAmount string `json:"subtotal_amount"`
	DiscountAmount string `json:"discount_amount"`
	ShippingFee    string `json:"shipping_fee"`
	TotalAmount    string `json:"total_amount"`
	PlacedAt       string `json:"placed_at"`
}

type PaginationDTO struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int   `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type RevenueOrdersResponseDTO struct {
	Items      []RevenueOrderItemDTO `json:"items"`
	Pagination PaginationDTO         `json:"pagination"`
}