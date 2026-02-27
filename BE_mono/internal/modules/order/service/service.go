package service

import (
	"database/sql"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"

	apperrors "golf-store/be-mono/internal/shared/errors"
	"golf-store/be-mono/internal/shared/model"
)

type Service struct {
	db *sql.DB
}

func New(db *sql.DB) *Service {
	return &Service{db: db}
}

type CreateOrderItemInput struct {
	ProductID string
	Quantity  int
}

type CreateOrderInput struct {
	UserID         string
	IdempotencyKey string
	Items          []CreateOrderItemInput
	Shipping       model.ShippingAddress
	CustomerNote   string
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
	Items      []*model.Order
	Page       int
	PageSize   int
	Total      int
	TotalPages int
}

func (s *Service) Create(input CreateOrderInput) (*model.Order, *apperrors.APIError) {
	if strings.TrimSpace(input.IdempotencyKey) == "" {
		return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: "Idempotency-Key is required"}
	}
	if len(input.Items) == 0 {
		return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "INVALID_ITEMS", Message: "order must include at least one item"}
	}

	// Idempotency check first.
	existingID := ""
	err := s.db.QueryRow(
		`SELECT id FROM orders WHERE user_id = ? AND idempotency_key = ? LIMIT 1`,
		input.UserID, strings.TrimSpace(input.IdempotencyKey),
	).Scan(&existingID)
	if err == nil {
		return s.GetByIDForUser(existingID, input.UserID)
	}
	if err != sql.ErrNoRows {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to check idempotency"}
	}

	now := time.Now().UTC()
	orderID := uuid.NewString()
	orderCode := fmt.Sprintf("ORD-%s-%s", now.Format("20060102"), strings.ToUpper(uuid.NewString()[:6]))

	tx, err := s.db.Begin()
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to begin transaction"}
	}
	defer tx.Rollback()

	var userStatus string
	if err := tx.QueryRow(`SELECT status FROM users WHERE id = ? LIMIT 1`, input.UserID).Scan(&userStatus); err != nil {
		if err == sql.ErrNoRows {
			return nil, &apperrors.APIError{Status: http.StatusUnauthorized, Code: "UNAUTHORIZED", Message: "Unauthorized"}
		}
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query user"}
	}
	if userStatus != string(model.UserStatusActive) {
		return nil, &apperrors.APIError{Status: http.StatusUnauthorized, Code: "UNAUTHORIZED", Message: "Unauthorized"}
	}

	orderItems := make([]model.OrderItem, 0, len(input.Items))
	var subtotal int64

	for _, item := range input.Items {
		if item.Quantity <= 0 {
			return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "INVALID_ITEMS", Message: "quantity must be greater than 0"}
		}

		var productID, sku, name, status string
		var priceCents int64
		var stock int
		var deletedAt sql.NullTime
		err := tx.QueryRow(
			`SELECT id, sku, name, price_cents, stock, status, deleted_at
			   FROM products
			  WHERE id = ?
			  FOR UPDATE`,
			item.ProductID,
		).Scan(&productID, &sku, &name, &priceCents, &stock, &status, &deletedAt)
		if err == sql.ErrNoRows {
			return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "INVALID_ITEMS", Message: "product not found"}
		}
		if err != nil {
			return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to lock product"}
		}
		if deletedAt.Valid || status != string(model.ProductStatusActive) {
			return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "INVALID_ITEMS", Message: "product not available"}
		}
		if stock < item.Quantity {
			return nil, &apperrors.APIError{Status: http.StatusConflict, Code: "OUT_OF_STOCK", Message: "Out of stock"}
		}

		if _, err := tx.Exec(
			`UPDATE products SET stock = stock - ?, updated_at = ? WHERE id = ?`,
			item.Quantity, now, productID,
		); err != nil {
			return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to update stock"}
		}

		lineTotal := priceCents * int64(item.Quantity)
		subtotal += lineTotal
		orderItems = append(orderItems, model.OrderItem{
			ID:        uuid.NewString(),
			OrderID:   orderID,
			ProductID: productID,
			SKU:       sku,
			Name:      name,
			UnitPrice: priceCents,
			Quantity:  item.Quantity,
			LineTotal: lineTotal,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	shippingFee := int64(3000000)
	discountAmount := int64(0)
	totalAmount := subtotal + shippingFee - discountAmount
	dueAt := now.Add(30 * time.Minute)

	if _, err := tx.Exec(
		`INSERT INTO orders (
			id, order_code, idempotency_key, user_id, order_status, payment_status, currency_code,
			subtotal_amount, discount_amount, shipping_fee, total_amount, payment_due_at,
			customer_note, cancel_reason,
			shipping_recipient_name, shipping_phone, shipping_line1, shipping_line2, shipping_ward,
			shipping_district, shipping_city, shipping_province, shipping_postal_code, shipping_country_code,
			placed_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		orderID, orderCode, strings.TrimSpace(input.IdempotencyKey), input.UserID,
		model.OrderStatusPendingPayment, model.PaymentStatusUnpaid, "VND",
		subtotal, discountAmount, shippingFee, totalAmount, dueAt,
		strings.TrimSpace(input.CustomerNote),
		input.Shipping.RecipientName,
		input.Shipping.RecipientPhone,
		input.Shipping.Line1,
		nullIfEmpty(input.Shipping.Line2),
		nullIfEmpty(input.Shipping.Ward),
		nullIfEmpty(input.Shipping.District),
		input.Shipping.City,
		nullIfEmpty(input.Shipping.Province),
		nullIfEmpty(input.Shipping.PostalCode),
		input.Shipping.CountryCode,
		now, now, now,
	); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to create order"}
	}

	for _, item := range orderItems {
		if _, err := tx.Exec(
			`INSERT INTO order_items (id, order_id, product_id, sku, name, unit_price, quantity, line_total, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			item.ID, item.OrderID, item.ProductID, item.SKU, item.Name, item.UnitPrice, item.Quantity, item.LineTotal, now, now,
		); err != nil {
			return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to create order item"}
		}
	}

	if err := addTrackingTx(tx, orderID, string(model.OrderStatusNew), string(model.OrderStatusPendingPayment), "SYSTEM", "Order created", now); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to create tracking event"}
	}
	if err := addNotificationTx(tx, input.UserID, "order.created", orderID, "Đơn hàng đã tạo", "Đơn hàng của bạn đã được tạo và đang chờ thanh toán", now); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to create notification"}
	}

	if err := tx.Commit(); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to commit order creation"}
	}

	return s.GetByIDForUser(orderID, input.UserID)
}

func (s *Service) List(input ListInput) ListOutput {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 20
	}

	where := []string{"1=1"}
	args := make([]any, 0)
	if !input.Admin {
		where = append(where, "user_id = ?")
		args = append(args, input.UserID)
	}
	if strings.TrimSpace(input.Status) != "" {
		where = append(where, "order_status = ?")
		args = append(args, strings.ToUpper(strings.TrimSpace(input.Status)))
	}
	if input.From != nil {
		where = append(where, "placed_at >= ?")
		args = append(args, *input.From)
	}
	if input.To != nil {
		where = append(where, "placed_at < ?")
		args = append(args, *input.To)
	}
	whereSQL := strings.Join(where, " AND ")

	var total int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM orders WHERE `+whereSQL, args...).Scan(&total); err != nil {
		return ListOutput{Items: []*model.Order{}, Page: input.Page, PageSize: input.PageSize, Total: 0, TotalPages: 0}
	}

	offset := (input.Page - 1) * input.PageSize
	listArgs := append(args, input.PageSize, offset)
	rows, err := s.db.Query(
		`SELECT id, order_code, user_id, order_status, payment_status, currency_code,
		        subtotal_amount, discount_amount, shipping_fee, total_amount, payment_due_at,
		        customer_note, cancel_reason,
		        shipping_recipient_name, shipping_phone, shipping_line1, shipping_line2, shipping_ward,
		        shipping_district, shipping_city, shipping_province, shipping_postal_code, shipping_country_code,
		        placed_at, created_at, updated_at
		   FROM orders
		  WHERE `+whereSQL+`
		  ORDER BY placed_at DESC
		  LIMIT ? OFFSET ?`,
		listArgs...,
	)
	if err != nil {
		return ListOutput{Items: []*model.Order{}, Page: input.Page, PageSize: input.PageSize, Total: total, TotalPages: calcTotalPages(total, input.PageSize)}
	}
	defer rows.Close()

	items := make([]*model.Order, 0)
	for rows.Next() {
		order, err := scanOrderSummary(rows)
		if err != nil {
			continue
		}
		items = append(items, order)
	}

	return ListOutput{
		Items:      items,
		Page:       input.Page,
		PageSize:   input.PageSize,
		Total:      total,
		TotalPages: calcTotalPages(total, input.PageSize),
	}
}

func (s *Service) GetByIDForUser(orderID string, userID string) (*model.Order, *apperrors.APIError) {
	order, err := s.getOrder(orderID)
	if err == sql.ErrNoRows {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query order"}
	}
	if order.UserID != userID {
		return nil, apperrors.ErrNotFound
	}
	return order, nil
}

func (s *Service) GetByIDAdmin(orderID string) (*model.Order, *apperrors.APIError) {
	order, err := s.getOrder(orderID)
	if err == sql.ErrNoRows {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query order"}
	}
	return order, nil
}

func (s *Service) CancelByUser(orderID string, userID string, reason string) (*model.Order, *apperrors.APIError) {
	now := time.Now().UTC()

	tx, err := s.db.Begin()
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to begin transaction"}
	}
	defer tx.Rollback()

	var status string
	err = tx.QueryRow(`SELECT order_status FROM orders WHERE id = ? AND user_id = ? FOR UPDATE`, orderID, userID).Scan(&status)
	if err == sql.ErrNoRows {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to lock order"}
	}

	current := model.OrderStatus(status)
	if !canCancel(current) {
		return nil, &apperrors.APIError{Status: http.StatusConflict, Code: "ORDER_CANNOT_CANCEL", Message: "Order cannot be canceled"}
	}

	if _, err := tx.Exec(
		`UPDATE orders SET order_status = ?, cancel_reason = ?, updated_at = ? WHERE id = ?`,
		model.OrderStatusCancelled, strings.TrimSpace(reason), now, orderID,
	); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to cancel order"}
	}

	rows, err := tx.Query(`SELECT product_id, quantity FROM order_items WHERE order_id = ?`, orderID)
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query order items"}
	}
	for rows.Next() {
		var productID string
		var qty int
		if err := rows.Scan(&productID, &qty); err != nil {
			continue
		}
		_, _ = tx.Exec(`UPDATE products SET stock = stock + ?, updated_at = ? WHERE id = ?`, qty, now, productID)
	}
	rows.Close()

	if err := addTrackingTx(tx, orderID, string(current), string(model.OrderStatusCancelled), "USER", "Order canceled by user", now); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to write tracking"}
	}
	if err := addNotificationTx(tx, userID, "order.cancelled", orderID, "Đơn hàng đã hủy", "Đơn hàng của bạn đã được hủy", now); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to create notification"}
	}

	if err := tx.Commit(); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to commit cancellation"}
	}

	return s.GetByIDForUser(orderID, userID)
}

func (s *Service) TrackingForUser(orderID string, userID string) ([]model.TrackingEvent, model.OrderStatus, *apperrors.APIError) {
	var currentStatus string
	if err := s.db.QueryRow(`SELECT order_status FROM orders WHERE id = ? AND user_id = ?`, orderID, userID).Scan(&currentStatus); err != nil {
		if err == sql.ErrNoRows {
			return nil, "", apperrors.ErrNotFound
		}
		return nil, "", &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query tracking"}
	}

	events, err := s.loadTracking(orderID)
	if err != nil {
		return nil, "", &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query tracking events"}
	}
	return events, model.OrderStatus(currentStatus), nil
}

func (s *Service) TrackingForAdmin(orderID string) ([]model.TrackingEvent, model.OrderStatus, *apperrors.APIError) {
	var currentStatus string
	if err := s.db.QueryRow(`SELECT order_status FROM orders WHERE id = ?`, orderID).Scan(&currentStatus); err != nil {
		if err == sql.ErrNoRows {
			return nil, "", apperrors.ErrNotFound
		}
		return nil, "", &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query tracking"}
	}

	events, err := s.loadTracking(orderID)
	if err != nil {
		return nil, "", &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query tracking events"}
	}
	return events, model.OrderStatus(currentStatus), nil
}

func (s *Service) AdminUpdateStatus(orderID string, toStatus model.OrderStatus, reason string) (*model.Order, *apperrors.APIError) {
	if toStatus == "" {
		return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: "to_status is required"}
	}
	now := time.Now().UTC()

	tx, err := s.db.Begin()
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to begin transaction"}
	}
	defer tx.Rollback()

	var userID string
	var current string
	if err := tx.QueryRow(`SELECT user_id, order_status FROM orders WHERE id = ? FOR UPDATE`, orderID).Scan(&userID, &current); err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to lock order"}
	}
	from := model.OrderStatus(current)
	if !isValidTransition(from, toStatus) {
		return nil, &apperrors.APIError{Status: http.StatusConflict, Code: "CONFLICT", Message: "invalid order status transition"}
	}

	if _, err := tx.Exec(`UPDATE orders SET order_status = ?, updated_at = ? WHERE id = ?`, toStatus, now, orderID); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to update order status"}
	}
	if toStatus == model.OrderStatusPaid {
		_, _ = tx.Exec(`UPDATE orders SET payment_status = ? WHERE id = ?`, model.PaymentStatusPaid, orderID)
	}

	if err := addTrackingTx(tx, orderID, string(from), string(toStatus), "ADMIN", strings.TrimSpace(reason), now); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to write tracking"}
	}
	if err := addNotificationTx(tx, userID, "order.status.changed", orderID, "Đơn hàng cập nhật trạng thái", fmt.Sprintf("Đơn hàng chuyển sang %s", toStatus), now); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to create notification"}
	}

	if err := tx.Commit(); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to commit status update"}
	}

	return s.GetByIDAdmin(orderID)
}

func (s *Service) MarkPaymentResult(orderID string, success bool, reason string) (*model.Order, *apperrors.APIError) {
	now := time.Now().UTC()
	tx, err := s.db.Begin()
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to begin transaction"}
	}
	defer tx.Rollback()

	var userID string
	var orderStatus string
	if err := tx.QueryRow(`SELECT user_id, order_status FROM orders WHERE id = ? FOR UPDATE`, orderID).Scan(&userID, &orderStatus); err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to lock order"}
	}

	if success {
		newStatus := model.OrderStatus(orderStatus)
		if newStatus == model.OrderStatusPendingPayment {
			newStatus = model.OrderStatusPaid
			if _, err := tx.Exec(`UPDATE orders SET order_status = ?, payment_status = ?, updated_at = ? WHERE id = ?`,
				newStatus, model.PaymentStatusPaid, now, orderID); err != nil {
				return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to update paid order"}
			}
			if err := addTrackingTx(tx, orderID, orderStatus, string(model.OrderStatusPaid), "PAYMENT_GATEWAY", "Payment succeeded", now); err != nil {
				return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to write tracking"}
			}
		} else {
			if _, err := tx.Exec(`UPDATE orders SET payment_status = ?, updated_at = ? WHERE id = ?`, model.PaymentStatusPaid, now, orderID); err != nil {
				return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to update payment status"}
			}
		}
		if err := addNotificationTx(tx, userID, "payment.succeeded", orderID, "Thanh toán thành công", "Thanh toán đơn hàng thành công", now); err != nil {
			return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to create notification"}
		}
	} else {
		if _, err := tx.Exec(`UPDATE orders SET payment_status = ?, updated_at = ? WHERE id = ?`, model.PaymentStatusFailed, now, orderID); err != nil {
			return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to update payment status"}
		}
		if err := addNotificationTx(tx, userID, "payment.failed", orderID, "Thanh toán thất bại", reason, now); err != nil {
			return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to create notification"}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to commit payment update"}
	}

	return s.GetByIDAdmin(orderID)
}

func (s *Service) loadTracking(orderID string) ([]model.TrackingEvent, error) {
	rows, err := s.db.Query(
		`SELECT id, order_id, from_status, to_status, source_type, description, occurred_at
		   FROM order_tracking_events
		  WHERE order_id = ?
		  ORDER BY occurred_at ASC`,
		orderID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]model.TrackingEvent, 0)
	for rows.Next() {
		var event model.TrackingEvent
		var fromStatus sql.NullString
		if err := rows.Scan(&event.ID, &event.OrderID, &fromStatus, &event.ToStatus, &event.SourceType, &event.Description, &event.OccurredAt); err != nil {
			continue
		}
		if fromStatus.Valid {
			event.FromStatus = fromStatus.String
		}
		events = append(events, event)
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].OccurredAt.Before(events[j].OccurredAt)
	})

	return events, nil
}

func (s *Service) getOrder(orderID string) (*model.Order, error) {
	row := s.db.QueryRow(
		`SELECT id, order_code, user_id, order_status, payment_status, currency_code,
		        subtotal_amount, discount_amount, shipping_fee, total_amount, payment_due_at,
		        customer_note, cancel_reason,
		        shipping_recipient_name, shipping_phone, shipping_line1, shipping_line2, shipping_ward,
		        shipping_district, shipping_city, shipping_province, shipping_postal_code, shipping_country_code,
		        placed_at, created_at, updated_at
		   FROM orders
		  WHERE id = ?`,
		orderID,
	)
	order, err := scanOrderSummary(row)
	if err != nil {
		return nil, err
	}

	itemRows, err := s.db.Query(
		`SELECT id, order_id, product_id, sku, name, unit_price, quantity, line_total, created_at, updated_at
		   FROM order_items
		  WHERE order_id = ?
		  ORDER BY created_at ASC`,
		orderID,
	)
	if err != nil {
		return nil, err
	}
	defer itemRows.Close()

	items := make([]model.OrderItem, 0)
	for itemRows.Next() {
		var item model.OrderItem
		if err := itemRows.Scan(
			&item.ID, &item.OrderID, &item.ProductID, &item.SKU, &item.Name,
			&item.UnitPrice, &item.Quantity, &item.LineTotal, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			continue
		}
		items = append(items, item)
	}
	order.Items = items

	return order, nil
}

func scanOrderSummary(scanner interface {
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

func addTrackingTx(tx *sql.Tx, orderID string, from string, to string, source string, description string, now time.Time) error {
	_, err := tx.Exec(
		`INSERT INTO order_tracking_events (id, order_id, from_status, to_status, source_type, description, occurred_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		uuid.NewString(), orderID, nullIfEmpty(from), to, source, nullIfEmpty(description), now,
	)
	return err
}

func addNotificationTx(tx *sql.Tx, userID string, eventType string, eventKey string, title string, content string, now time.Time) error {
	_, err := tx.Exec(
		`INSERT INTO notifications (id, user_id, channel, event_type, event_key, title, content, status, is_read, sent_at, created_at, updated_at)
		 VALUES (?, ?, 'IN_APP', ?, ?, ?, ?, 'SENT', 0, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE updated_at = VALUES(updated_at)`,
		uuid.NewString(), userID, eventType, eventKey, title, content, now, now, now,
	)
	return err
}

func nullIfEmpty(v string) any {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	return strings.TrimSpace(v)
}

func canCancel(status model.OrderStatus) bool {
	switch status {
	case model.OrderStatusNew, model.OrderStatusPendingPayment, model.OrderStatusPaid, model.OrderStatusProcessing:
		return true
	default:
		return false
	}
}

func isValidTransition(from model.OrderStatus, to model.OrderStatus) bool {
	if from == to {
		return true
	}
	allowed := map[model.OrderStatus][]model.OrderStatus{
		model.OrderStatusNew:            {model.OrderStatusPendingPayment, model.OrderStatusCancelled},
		model.OrderStatusPendingPayment: {model.OrderStatusPaid, model.OrderStatusCancelled},
		model.OrderStatusPaid:           {model.OrderStatusProcessing, model.OrderStatusCancelled},
		model.OrderStatusProcessing:     {model.OrderStatusShipping, model.OrderStatusCancelled},
		model.OrderStatusShipping:       {model.OrderStatusCompleted},
		model.OrderStatusCompleted:      {},
		model.OrderStatusCancelled:      {},
	}
	for _, candidate := range allowed[from] {
		if candidate == to {
			return true
		}
	}
	return false
}

func ParseOrderStatus(value string) (model.OrderStatus, error) {
	up := strings.ToUpper(strings.TrimSpace(value))
	switch model.OrderStatus(up) {
	case model.OrderStatusNew,
		model.OrderStatusPendingPayment,
		model.OrderStatusPaid,
		model.OrderStatusProcessing,
		model.OrderStatusShipping,
		model.OrderStatusCompleted,
		model.OrderStatusCancelled:
		return model.OrderStatus(up), nil
	default:
		return "", fmt.Errorf("invalid order status")
	}
}

func calcTotalPages(total int, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	return (total + pageSize - 1) / pageSize
}
