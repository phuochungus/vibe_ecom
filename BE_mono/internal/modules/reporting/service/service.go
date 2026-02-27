package service

import (
	"database/sql"
	"fmt"
	"time"

	"golf-store/be-mono/internal/shared/model"
)

type Service struct {
	db *sql.DB
}

func New(db *sql.DB) *Service {
	return &Service{db: db}
}

type SummaryOutput struct {
	From               time.Time
	To                 time.Time
	GrossRevenue       int64
	RefundAmount       int64
	NetRevenue         int64
	CompletedOrders    int
	PaymentSuccessRate string
}

type OrdersOutput struct {
	Items      []*model.Order
	Page       int
	PageSize   int
	Total      int
	TotalPages int
}

func (s *Service) Summary(from time.Time, to time.Time) SummaryOutput {
	var grossRevenue int64
	var completedOrders int
	_ = s.db.QueryRow(
		`SELECT COALESCE(SUM(total_amount), 0), COUNT(*)
		   FROM orders
		  WHERE order_status = 'COMPLETED' AND placed_at >= ? AND placed_at < ?`,
		from, to,
	).Scan(&grossRevenue, &completedOrders)

	var refundAmount int64
	_ = s.db.QueryRow(
		`SELECT COALESCE(SUM(amount), 0)
		   FROM payment_transactions
		  WHERE txn_type = 'REFUND' AND status = 'SUCCESS' AND created_at >= ? AND created_at < ?`,
		from, to,
	).Scan(&refundAmount)

	var paymentTotal int
	var paymentSuccess int
	_ = s.db.QueryRow(
		`SELECT COUNT(*), COALESCE(SUM(CASE WHEN status = 'SUCCESS' THEN 1 ELSE 0 END), 0)
		   FROM payment_transactions
		  WHERE txn_type = 'PAYMENT' AND created_at >= ? AND created_at < ?`,
		from, to,
	).Scan(&paymentTotal, &paymentSuccess)

	rate := "0.00"
	if paymentTotal > 0 {
		rate = fmt.Sprintf("%.2f", float64(paymentSuccess)/float64(paymentTotal))
	}

	return SummaryOutput{
		From:               from,
		To:                 to,
		GrossRevenue:       grossRevenue,
		RefundAmount:       refundAmount,
		NetRevenue:         grossRevenue - refundAmount,
		CompletedOrders:    completedOrders,
		PaymentSuccessRate: rate,
	}
}

func (s *Service) Orders(from time.Time, to time.Time, page int, pageSize int) OrdersOutput {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	var total int
	_ = s.db.QueryRow(
		`SELECT COUNT(*)
		   FROM orders
		  WHERE order_status = 'COMPLETED' AND placed_at >= ? AND placed_at < ?`,
		from, to,
	).Scan(&total)

	offset := (page - 1) * pageSize
	rows, err := s.db.Query(
		`SELECT id, order_code, user_id, order_status, payment_status, currency_code,
		        subtotal_amount, discount_amount, shipping_fee, total_amount, payment_due_at,
		        customer_note, cancel_reason,
		        shipping_recipient_name, shipping_phone, shipping_line1, shipping_line2, shipping_ward,
		        shipping_district, shipping_city, shipping_province, shipping_postal_code, shipping_country_code,
		        placed_at, created_at, updated_at
		   FROM orders
		  WHERE order_status = 'COMPLETED' AND placed_at >= ? AND placed_at < ?
		  ORDER BY placed_at DESC
		  LIMIT ? OFFSET ?`,
		from, to, pageSize, offset,
	)
	if err != nil {
		return OrdersOutput{Items: []*model.Order{}, Page: page, PageSize: pageSize, Total: total, TotalPages: calcTotalPages(total, pageSize)}
	}
	defer rows.Close()

	items := make([]*model.Order, 0)
	for rows.Next() {
		order, err := scanOrderForReport(rows)
		if err != nil {
			continue
		}
		items = append(items, order)
	}

	return OrdersOutput{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: calcTotalPages(total, pageSize),
	}
}

func scanOrderForReport(scanner interface {
	Scan(dest ...any) error
}) (*model.Order, error) {
	order := &model.Order{}
	var orderStatus, paymentStatus string
	var paymentDueAt sql.NullTime
	var customerNote, cancelReason sql.NullString
	var line2, ward, district, province, postalCode sql.NullString

	err := scanner.Scan(
		&order.ID,
		&order.OrderCode,
		&order.UserID,
		&orderStatus,
		&paymentStatus,
		&order.CurrencyCode,
		&order.SubtotalAmount,
		&order.DiscountAmount,
		&order.ShippingFee,
		&order.TotalAmount,
		&paymentDueAt,
		&customerNote,
		&cancelReason,
		&order.Shipping.RecipientName,
		&order.Shipping.RecipientPhone,
		&order.Shipping.Line1,
		&line2,
		&ward,
		&district,
		&order.Shipping.City,
		&province,
		&postalCode,
		&order.Shipping.CountryCode,
		&order.PlacedAt,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	order.OrderStatus = model.OrderStatus(orderStatus)
	order.PaymentStatus = model.PaymentStatus(paymentStatus)
	if paymentDueAt.Valid {
		t := paymentDueAt.Time.UTC()
		order.PaymentDueAt = &t
	}
	if customerNote.Valid {
		order.CustomerNote = customerNote.String
	}
	if cancelReason.Valid {
		order.CancelReason = cancelReason.String
	}
	if line2.Valid {
		order.Shipping.Line2 = line2.String
	}
	if ward.Valid {
		order.Shipping.Ward = ward.String
	}
	if district.Valid {
		order.Shipping.District = district.String
	}
	if province.Valid {
		order.Shipping.Province = province.String
	}
	if postalCode.Valid {
		order.Shipping.PostalCode = postalCode.String
	}

	return order, nil
}

func calcTotalPages(total int, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	return (total + pageSize - 1) / pageSize
}
