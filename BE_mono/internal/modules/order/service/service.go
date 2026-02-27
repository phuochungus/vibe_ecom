package service

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"golf-store/be-mono/internal/platform/db"
	apperrors "golf-store/be-mono/internal/shared/errors"
)

type Service struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Service {
	return &Service{db: db}
}

type CreateOrderItemInput struct {
	ProductID string
	Quantity  int
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
	Items      []*db.OrderEntity
	Page       int
	PageSize   int
	Total      int
	TotalPages int
}

func (s *Service) Create(input CreateOrderInput) (*db.OrderEntity, *apperrors.APIError) {
	idemKey := strings.TrimSpace(input.IdempotencyKey)
	if idemKey == "" {
		return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: "Idempotency-Key is required"}
	}
	if len(input.Items) == 0 {
		return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "INVALID_ITEMS", Message: "order must include at least one item"}
	}

	var existing db.OrderEntity
	err := s.db.Where("user_id = ? AND idempotency_key = ?", input.UserID, idemKey).Take(&existing).Error
	if err == nil {
		return s.GetByIDForUser(existing.ID, input.UserID)
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to check idempotency"}
	}

	now := time.Now().UTC()
	orderID := uuid.NewString()
	orderCode := fmt.Sprintf("ORD-%s-%s", now.Format("20060102"), strings.ToUpper(uuid.NewString()[:6]))

	err = s.db.Transaction(func(tx *gorm.DB) error {
		var user db.UserEntity
		if err := tx.Select("id", "status").Where("id = ?", input.UserID).Take(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &apperrors.APIError{Status: http.StatusUnauthorized, Code: "UNAUTHORIZED", Message: "Unauthorized"}
			}
			return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query user"}
		}
		if user.Status != db.UserStatusActive {
			return &apperrors.APIError{Status: http.StatusUnauthorized, Code: "UNAUTHORIZED", Message: "Unauthorized"}
		}

		orderItems := make([]db.OrderItemEntity, 0, len(input.Items))
		var subtotal int64

		for _, item := range input.Items {
			if item.Quantity <= 0 {
				return &apperrors.APIError{Status: http.StatusBadRequest, Code: "INVALID_ITEMS", Message: "quantity must be greater than 0"}
			}

			var product db.ProductEntity
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", item.ProductID).Take(&product).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return &apperrors.APIError{Status: http.StatusBadRequest, Code: "INVALID_ITEMS", Message: "product not found"}
				}
				return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to lock product"}
			}

			if product.DeletedAt != nil || product.Status != db.ProductStatusActive {
				return &apperrors.APIError{Status: http.StatusBadRequest, Code: "INVALID_ITEMS", Message: "product not available"}
			}
			if product.Stock < item.Quantity {
				return &apperrors.APIError{Status: http.StatusConflict, Code: "OUT_OF_STOCK", Message: "Out of stock"}
			}

			if err := tx.Model(&db.ProductEntity{}).
				Where("id = ?", product.ID).
				Updates(map[string]any{
					"stock":      gorm.Expr("stock - ?", item.Quantity),
					"updated_at": now,
				}).Error; err != nil {
				return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to update stock"}
			}

			lineTotal := product.Price * int64(item.Quantity)
			subtotal += lineTotal
			orderItems = append(orderItems, db.OrderItemEntity{
				ID:        uuid.NewString(),
				OrderID:   orderID,
				ProductID: product.ID,
				SKU:       product.SKU,
				Name:      product.Name,
				UnitPrice: product.Price,
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

		order := db.OrderEntity{
			ID:                    orderID,
			OrderCode:             orderCode,
			IdempotencyKey:        idemKey,
			UserID:                input.UserID,
			OrderStatus:           db.OrderStatusPendingPayment,
			PaymentStatus:         db.PaymentStatusUnpaid,
			CurrencyCode:          "VND",
			SubtotalAmount:        subtotal,
			DiscountAmount:        discountAmount,
			ShippingFee:           shippingFee,
			TotalAmount:           totalAmount,
			PaymentDueAt:          &dueAt,
			CustomerNote:          stringPtrOrNil(input.CustomerNote),
			ShippingRecipientName: strings.TrimSpace(input.Shipping.RecipientName),
			ShippingPhone:         strings.TrimSpace(input.Shipping.RecipientPhone),
			ShippingLine1:         strings.TrimSpace(input.Shipping.Line1),
			ShippingLine2:         stringPtrOrNil(input.Shipping.Line2),
			ShippingWard:          stringPtrOrNil(input.Shipping.Ward),
			ShippingDistrict:      stringPtrOrNil(input.Shipping.District),
			ShippingCity:          strings.TrimSpace(input.Shipping.City),
			ShippingProvince:      stringPtrOrNil(input.Shipping.Province),
			ShippingPostalCode:    stringPtrOrNil(input.Shipping.PostalCode),
			ShippingCountryCode:   strings.TrimSpace(input.Shipping.CountryCode),
			PlacedAt:              now,
			CreatedAt:             now,
			UpdatedAt:             now,
		}
		if err := tx.Create(&order).Error; err != nil {
			return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to create order"}
		}

		if len(orderItems) > 0 {
			if err := tx.Create(&orderItems).Error; err != nil {
				return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to create order item"}
			}
		}

		if err := addTrackingTx(tx, orderID, db.OrderStatusNew, db.OrderStatusPendingPayment, "SYSTEM", "Order created", now); err != nil {
			return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to create tracking event"}
		}
		if err := addNotificationTx(tx, input.UserID, "order.created", orderID, "Đơn hàng đã tạo", "Đơn hàng của bạn đã được tạo và đang chờ thanh toán", now); err != nil {
			return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to create notification"}
		}

		return nil
	})
	if err != nil {
		if apiErr, ok := err.(*apperrors.APIError); ok {
			return nil, apiErr
		}
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

	query := s.db.Model(&db.OrderEntity{})
	if !input.Admin {
		query = query.Where("user_id = ?", input.UserID)
	}
	if strings.TrimSpace(input.Status) != "" {
		query = query.Where("order_status = ?", strings.ToUpper(strings.TrimSpace(input.Status)))
	}
	if input.From != nil {
		query = query.Where("placed_at >= ?", *input.From)
	}
	if input.To != nil {
		query = query.Where("placed_at < ?", *input.To)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return ListOutput{Items: []*db.OrderEntity{}, Page: input.Page, PageSize: input.PageSize, Total: 0, TotalPages: 0}
	}

	offset := (input.Page - 1) * input.PageSize
	entities := make([]db.OrderEntity, 0)
	if err := query.Order("placed_at DESC").Limit(input.PageSize).Offset(offset).Find(&entities).Error; err != nil {
		return ListOutput{
			Items:      []*db.OrderEntity{},
			Page:       input.Page,
			PageSize:   input.PageSize,
			Total:      int(total),
			TotalPages: calcTotalPages(int(total), input.PageSize),
		}
	}

	items := make([]*db.OrderEntity, 0, len(entities))
	for i := range entities {
		items = append(items, cloneOrder(&entities[i]))
	}

	return ListOutput{
		Items:      items,
		Page:       input.Page,
		PageSize:   input.PageSize,
		Total:      int(total),
		TotalPages: calcTotalPages(int(total), input.PageSize),
	}
}

func (s *Service) GetByIDForUser(orderID string, userID string) (*db.OrderEntity, *apperrors.APIError) {
	order, err := s.getOrder(orderID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
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

func (s *Service) GetByIDAdmin(orderID string) (*db.OrderEntity, *apperrors.APIError) {
	order, err := s.getOrder(orderID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query order"}
	}
	return order, nil
}

func (s *Service) CancelByUser(orderID string, userID string, reason string) (*db.OrderEntity, *apperrors.APIError) {
	now := time.Now().UTC()

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var orderEntity db.OrderEntity
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND user_id = ?", orderID, userID).
			Take(&orderEntity).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperrors.ErrNotFound
			}
			return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to lock order"}
		}

		current := orderEntity.OrderStatus
		if !canCancel(current) {
			return &apperrors.APIError{Status: http.StatusConflict, Code: "ORDER_CANNOT_CANCEL", Message: "Order cannot be canceled"}
		}

		updates := map[string]any{
			"order_status": db.OrderStatusCancelled,
			"updated_at":   now,
		}
		if trimmedReason := strings.TrimSpace(reason); trimmedReason != "" {
			updates["cancel_reason"] = trimmedReason
		} else {
			updates["cancel_reason"] = nil
		}
		if err := tx.Model(&db.OrderEntity{}).Where("id = ?", orderID).Updates(updates).Error; err != nil {
			return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to cancel order"}
		}

		items := make([]db.OrderItemEntity, 0)
		if err := tx.Where("order_id = ?", orderID).Find(&items).Error; err != nil {
			return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query order items"}
		}
		for _, item := range items {
			if err := tx.Model(&db.ProductEntity{}).
				Where("id = ?", item.ProductID).
				Updates(map[string]any{
					"stock":      gorm.Expr("stock + ?", item.Quantity),
					"updated_at": now,
				}).Error; err != nil {
				return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to restore stock"}
			}
		}

		if err := addTrackingTx(tx, orderID, current, db.OrderStatusCancelled, "USER", "Order canceled by user", now); err != nil {
			return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to write tracking"}
		}
		if err := addNotificationTx(tx, userID, "order.cancelled", orderID, "Đơn hàng đã hủy", "Đơn hàng của bạn đã được hủy", now); err != nil {
			return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to create notification"}
		}

		return nil
	})
	if err != nil {
		if apiErr, ok := err.(*apperrors.APIError); ok {
			return nil, apiErr
		}
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to commit cancellation"}
	}

	return s.GetByIDForUser(orderID, userID)
}

func (s *Service) TrackingForUser(orderID string, userID string) ([]db.OrderTrackingEventEntity, string, *apperrors.APIError) {
	var orderEntity db.OrderEntity
	if err := s.db.Select("id", "order_status").Where("id = ? AND user_id = ?", orderID, userID).Take(&orderEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", apperrors.ErrNotFound
		}
		return nil, "", &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query tracking"}
	}

	events, err := s.loadTracking(orderID)
	if err != nil {
		return nil, "", &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query tracking events"}
	}
	return events, orderEntity.OrderStatus, nil
}

func (s *Service) TrackingForAdmin(orderID string) ([]db.OrderTrackingEventEntity, string, *apperrors.APIError) {
	var orderEntity db.OrderEntity
	if err := s.db.Select("id", "order_status").Where("id = ?", orderID).Take(&orderEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", apperrors.ErrNotFound
		}
		return nil, "", &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query tracking"}
	}

	events, err := s.loadTracking(orderID)
	if err != nil {
		return nil, "", &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query tracking events"}
	}
	return events, orderEntity.OrderStatus, nil
}

func (s *Service) AdminUpdateStatus(orderID string, toStatus string, reason string) (*db.OrderEntity, *apperrors.APIError) {
	if toStatus == "" {
		return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: "to_status is required"}
	}
	now := time.Now().UTC()

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var orderEntity db.OrderEntity
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", orderID).
			Take(&orderEntity).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperrors.ErrNotFound
			}
			return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to lock order"}
		}

		from := orderEntity.OrderStatus
		if !isValidTransition(from, toStatus) {
			return &apperrors.APIError{Status: http.StatusConflict, Code: "CONFLICT", Message: "invalid order status transition"}
		}

		updates := map[string]any{
			"order_status": toStatus,
			"updated_at":   now,
		}
		if toStatus == db.OrderStatusPaid {
			updates["payment_status"] = db.PaymentStatusPaid
		}
		if err := tx.Model(&db.OrderEntity{}).Where("id = ?", orderID).Updates(updates).Error; err != nil {
			return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to update order status"}
		}

		if err := addTrackingTx(tx, orderID, from, toStatus, "ADMIN", strings.TrimSpace(reason), now); err != nil {
			return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to write tracking"}
		}
		if err := addNotificationTx(tx, orderEntity.UserID, "order.status.changed", orderID, "Đơn hàng cập nhật trạng thái", fmt.Sprintf("Đơn hàng chuyển sang %s", toStatus), now); err != nil {
			return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to create notification"}
		}

		return nil
	})
	if err != nil {
		if apiErr, ok := err.(*apperrors.APIError); ok {
			return nil, apiErr
		}
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to commit status update"}
	}

	return s.GetByIDAdmin(orderID)
}

func (s *Service) MarkPaymentResult(orderID string, success bool, reason string) (*db.OrderEntity, *apperrors.APIError) {
	now := time.Now().UTC()
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var orderEntity db.OrderEntity
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", orderID).
			Take(&orderEntity).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperrors.ErrNotFound
			}
			return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to lock order"}
		}

		if success {
			if orderEntity.OrderStatus == db.OrderStatusPendingPayment {
				if err := tx.Model(&db.OrderEntity{}).
					Where("id = ?", orderID).
					Updates(map[string]any{
						"order_status":   db.OrderStatusPaid,
						"payment_status": db.PaymentStatusPaid,
						"updated_at":     now,
					}).Error; err != nil {
					return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to update paid order"}
				}
				if err := addTrackingTx(tx, orderID, orderEntity.OrderStatus, db.OrderStatusPaid, "PAYMENT_GATEWAY", "Payment succeeded", now); err != nil {
					return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to write tracking"}
				}
			} else {
				if err := tx.Model(&db.OrderEntity{}).
					Where("id = ?", orderID).
					Updates(map[string]any{
						"payment_status": db.PaymentStatusPaid,
						"updated_at":     now,
					}).Error; err != nil {
					return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to update payment status"}
				}
			}
			if err := addNotificationTx(tx, orderEntity.UserID, "payment.succeeded", orderID, "Thanh toán thành công", "Thanh toán đơn hàng thành công", now); err != nil {
				return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to create notification"}
			}
		} else {
			if err := tx.Model(&db.OrderEntity{}).
				Where("id = ?", orderID).
				Updates(map[string]any{
					"payment_status": db.PaymentStatusFailed,
					"updated_at":     now,
				}).Error; err != nil {
				return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to update payment status"}
			}
			if err := addNotificationTx(tx, orderEntity.UserID, "payment.failed", orderID, "Thanh toán thất bại", reason, now); err != nil {
				return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to create notification"}
			}
		}

		return nil
	})
	if err != nil {
		if apiErr, ok := err.(*apperrors.APIError); ok {
			return nil, apiErr
		}
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to commit payment update"}
	}

	return s.GetByIDAdmin(orderID)
}

func (s *Service) loadTracking(orderID string) ([]db.OrderTrackingEventEntity, error) {
	entities := make([]db.OrderTrackingEventEntity, 0)
	if err := s.db.Where("order_id = ?", orderID).Order("occurred_at ASC").Find(&entities).Error; err != nil {
		return nil, err
	}

	events := make([]db.OrderTrackingEventEntity, 0, len(entities))
	for i := range entities {
		events = append(events, *cloneTracking(&entities[i]))
	}

	return events, nil
}

func (s *Service) getOrder(orderID string) (*db.OrderEntity, error) {
	var orderEntity db.OrderEntity
	if err := s.db.
		Preload("Items").
		Where("id = ?", orderID).
		Take(&orderEntity).Error; err != nil {
		return nil, err
	}
	return &orderEntity, nil
}

func addTrackingTx(tx *gorm.DB, orderID string, from string, to string, source string, description string, now time.Time) error {
	event := db.OrderTrackingEventEntity{
		ID:          uuid.NewString(),
		OrderID:     orderID,
		FromStatus:  stringPtrOrNil(from),
		ToStatus:    strings.TrimSpace(to),
		SourceType:  strings.TrimSpace(source),
		Description: stringPtrOrNil(description),
		OccurredAt:  now,
	}
	return tx.Create(&event).Error
}

func addNotificationTx(tx *gorm.DB, userID string, eventType string, eventKey string, title string, content string, now time.Time) error {
	notification := db.NotificationEntity{
		ID:        uuid.NewString(),
		UserID:    userID,
		Channel:   "IN_APP",
		EventType: strings.TrimSpace(eventType),
		EventKey:  strings.TrimSpace(eventKey),
		Title:     strings.TrimSpace(title),
		Content:   strings.TrimSpace(content),
		Status:    db.NotificationStatusSent,
		IsRead:    false,
		SentAt:    &now,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "event_type"},
			{Name: "event_key"},
		},
		DoUpdates: clause.Assignments(map[string]any{"updated_at": now}),
	}).Create(&notification).Error
}

func stringPtrOrNil(v string) *string {
	trimmed := strings.TrimSpace(v)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func canCancel(status string) bool {
	switch status {
	case db.OrderStatusNew, db.OrderStatusPendingPayment, db.OrderStatusPaid, db.OrderStatusProcessing:
		return true
	default:
		return false
	}
}

func isValidTransition(from string, to string) bool {
	if from == to {
		return true
	}
	allowed := map[string][]string{
		db.OrderStatusNew:            {db.OrderStatusPendingPayment, db.OrderStatusCancelled},
		db.OrderStatusPendingPayment: {db.OrderStatusPaid, db.OrderStatusCancelled},
		db.OrderStatusPaid:           {db.OrderStatusProcessing, db.OrderStatusCancelled},
		db.OrderStatusProcessing:     {db.OrderStatusShipping, db.OrderStatusCancelled},
		db.OrderStatusShipping:       {db.OrderStatusCompleted},
		db.OrderStatusCompleted:      {},
		db.OrderStatusCancelled:      {},
	}
	for _, candidate := range allowed[from] {
		if candidate == to {
			return true
		}
	}
	return false
}

func ParseOrderStatus(value string) (string, error) {
	up := strings.ToUpper(strings.TrimSpace(value))
	switch up {
	case db.OrderStatusNew,
		db.OrderStatusPendingPayment,
		db.OrderStatusPaid,
		db.OrderStatusProcessing,
		db.OrderStatusShipping,
		db.OrderStatusCompleted,
		db.OrderStatusCancelled:
		return up, nil
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

func cloneOrder(entity *db.OrderEntity) *db.OrderEntity {
	if entity == nil {
		return nil
	}
	copy := *entity
	copy.CreatedAt = copy.CreatedAt.UTC()
	copy.UpdatedAt = copy.UpdatedAt.UTC()
	copy.PlacedAt = copy.PlacedAt.UTC()
	if copy.PaymentDueAt != nil {
		t := copy.PaymentDueAt.UTC()
		copy.PaymentDueAt = &t
	}
	return &copy
}

func cloneOrderItem(entity *db.OrderItemEntity) *db.OrderItemEntity {
	if entity == nil {
		return nil
	}
	copy := *entity
	copy.CreatedAt = copy.CreatedAt.UTC()
	copy.UpdatedAt = copy.UpdatedAt.UTC()
	return &copy
}

func cloneTracking(entity *db.OrderTrackingEventEntity) *db.OrderTrackingEventEntity {
	if entity == nil {
		return nil
	}
	copy := *entity
	copy.OccurredAt = copy.OccurredAt.UTC()
	return &copy
}
