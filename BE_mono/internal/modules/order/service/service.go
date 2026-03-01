package service

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"golf-store/be-mono/internal/modules/order/repository"
	entities "golf-store/be-mono/internal/platform/entities"
	apperrors "golf-store/be-mono/internal/shared/errors"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
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
	Items      []*entities.Order
	Page       int
	PageSize   int
	Total      int
	TotalPages int
}

func (s *Service) Create(input CreateOrderInput) (*entities.Order, *apperrors.APIError) {
	idemKey := strings.TrimSpace(input.IdempotencyKey)
	if idemKey == "" {
		return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: "Idempotency-Key is required"}
	}
	if len(input.Items) == 0 {
		return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "INVALID_ITEMS", Message: "order must include at least one item"}
	}

	existing, err := s.repo.FindIdempotentOrder(input.UserID, idemKey)
	if err == nil && existing != nil && existing.ID != "" {
		return s.GetByIDForUser(existing.ID, input.UserID)
	}
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to check idempotency"}
	}

	now := time.Now().UTC()
	orderID := uuid.NewString()
	orderCode := fmt.Sprintf("ORD-%s-%s", now.Format("20060102"), strings.ToUpper(uuid.NewString()[:6]))

	user, err := s.repo.FindUserByID(input.UserID)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, &apperrors.APIError{Status: http.StatusUnauthorized, Code: "UNAUTHORIZED", Message: "Unauthorized"}
		}
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query user"}
	}
	if user.Status != entities.UserStatusActive {
		return nil, &apperrors.APIError{Status: http.StatusUnauthorized, Code: "UNAUTHORIZED", Message: "Unauthorized"}
	}

	orderItems := make([]entities.OrderItem, 0, len(input.Items))
	var subtotal int64
	productUpdates := make(map[string]int)

	for _, item := range input.Items {
		if item.Quantity <= 0 {
			return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "INVALID_ITEMS", Message: "quantity must be greater than 0"}
		}

		product, err := s.repo.LockProduct(nil, item.ProductID)
		if err != nil {
			if strings.Contains(err.Error(), "record not found") {
				return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "INVALID_ITEMS", Message: "product not found"}
			}
			return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to lock product"}
		}

		if product.DeletedAt != nil || product.Status != entities.ProductStatusActive {
			return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "INVALID_ITEMS", Message: "product not available"}
		}
		if product.Stock < item.Quantity {
			return nil, &apperrors.APIError{Status: http.StatusConflict, Code: "OUT_OF_STOCK", Message: "Out of stock"}
		}

		productUpdates[product.ID] = item.Quantity

		lineTotal := product.Price * int64(item.Quantity)
		subtotal += lineTotal
		orderItems = append(orderItems, entities.OrderItem{
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

	order := entities.Order{
		ID:                    orderID,
		OrderCode:             orderCode,
		IdempotencyKey:        idemKey,
		UserID:                input.UserID,
		OrderStatus:           entities.OrderStatusPendingPayment,
		PaymentStatus:         entities.PaymentStatusUnpaid,
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

	tracking := &entities.OrderTrackingEvent{
		ID:          uuid.NewString(),
		OrderID:     orderID,
		FromStatus:  stringPtrOrNil(entities.OrderStatusNew),
		ToStatus:    entities.OrderStatusPendingPayment,
		SourceType:  "SYSTEM",
		Description: stringPtrOrNil("Order created"),
		OccurredAt:  now,
	}

	notification := &entities.Notification{
		ID:        uuid.NewString(),
		UserID:    input.UserID,
		Channel:   "IN_APP",
		EventType: "order.created",
		EventKey:  orderID,
		Title:     "Đơn hàng đã tạo",
		Content:   "Đơn hàng của bạn đã được tạo và đang chờ thanh toán",
		Status:    entities.NotificationStatusSent,
		IsRead:    false,
		SentAt:    &now,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.repo.CreateOrderTx(&order, orderItems, tracking, notification, productUpdates)
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

	filter := repository.ListFilter{
		UserID:   input.UserID,
		Status:   input.Status,
		From:     input.From,
		To:       input.To,
		Page:     input.Page,
		PageSize: input.PageSize,
		Admin:    input.Admin,
	}

	rows, total, err := s.repo.List(filter)
	if err != nil {
		return ListOutput{
			Items:      []*entities.Order{},
			Page:       input.Page,
			PageSize:   input.PageSize,
			Total:      0,
			TotalPages: 0,
		}
	}

	items := make([]*entities.Order, 0, len(rows))
	for i := range rows {
		items = append(items, cloneOrder(&rows[i]))
	}

	return ListOutput{
		Items:      items,
		Page:       input.Page,
		PageSize:   input.PageSize,
		Total:      int(total),
		TotalPages: calcTotalPages(int(total), input.PageSize),
	}
}

func (s *Service) GetByIDForUser(orderID string, userID string) (*entities.Order, *apperrors.APIError) {
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

func (s *Service) GetByIDAdmin(orderID string) (*entities.Order, *apperrors.APIError) {
	order, err := s.getOrder(orderID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query order"}
	}
	return order, nil
}

func (s *Service) CancelByUser(orderID string, userID string, reason string) (*entities.Order, *apperrors.APIError) {
	now := time.Now().UTC()

	updates := map[string]any{
		"order_status": entities.OrderStatusCancelled,
		"updated_at":   now,
	}
	if trimmedReason := strings.TrimSpace(reason); trimmedReason != "" {
		updates["cancel_reason"] = trimmedReason
	} else {
		updates["cancel_reason"] = nil
	}

	tracking := &entities.OrderTrackingEvent{
		ID:          uuid.NewString(),
		OrderID:     orderID,
		ToStatus:    entities.OrderStatusCancelled,
		SourceType:  "USER",
		Description: stringPtrOrNil("Order canceled by user"),
		OccurredAt:  now,
	}

	notification := &entities.Notification{
		ID:        uuid.NewString(),
		UserID:    userID,
		Channel:   "IN_APP",
		EventType: "order.cancelled",
		EventKey:  orderID,
		Title:     "Đơn hàng đã hủy",
		Content:   "Đơn hàng của bạn đã được hủy",
		Status:    entities.NotificationStatusSent,
		IsRead:    false,
		SentAt:    &now,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := s.repo.CancelOrderTx(orderID, userID, updates, tracking, notification)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, apperrors.ErrNotFound
		}
		if strings.Contains(err.Error(), "ORDER_CANNOT_CANCEL") {
			return nil, &apperrors.APIError{Status: http.StatusConflict, Code: "ORDER_CANNOT_CANCEL", Message: "Order cannot be canceled"}
		}
		if apiErr, ok := err.(*apperrors.APIError); ok {
			return nil, apiErr
		}
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to commit cancellation"}
	}

	return s.GetByIDForUser(orderID, userID)
}

func (s *Service) TrackingForUser(orderID string, userID string) ([]entities.OrderTrackingEvent, string, *apperrors.APIError) {
	order, err := s.repo.FindByID(orderID)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, "", apperrors.ErrNotFound
		}
		return nil, "", &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query tracking"}
	}
	if order.UserID != userID {
		return nil, "", apperrors.ErrNotFound
	}

	events, err := s.repo.LoadTracking(orderID)
	if err != nil {
		return nil, "", &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query tracking events"}
	}
	return events, order.OrderStatus, nil
}

func (s *Service) TrackingForAdmin(orderID string) ([]entities.OrderTrackingEvent, string, *apperrors.APIError) {
	order, err := s.repo.FindByID(orderID)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, "", apperrors.ErrNotFound
		}
		return nil, "", &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query tracking"}
	}

	events, err := s.repo.LoadTracking(orderID)
	if err != nil {
		return nil, "", &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query tracking events"}
	}
	return events, order.OrderStatus, nil
}

func (s *Service) AdminUpdateStatus(orderID string, toStatus string, reason string) (*entities.Order, *apperrors.APIError) {
	if toStatus == "" {
		return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: "to_status is required"}
	}
	now := time.Now().UTC()

	// get the order first for user ID (for notification) and current status
	order, err := s.repo.FindByID(orderID)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, apperrors.ErrNotFound
		}
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query order"}
	}

	updates := map[string]any{
		"order_status": toStatus,
		"updated_at":   now,
	}
	if toStatus == entities.OrderStatusPaid {
		updates["payment_status"] = entities.PaymentStatusPaid
	}

	tracking := &entities.OrderTrackingEvent{
		ID:          uuid.NewString(),
		OrderID:     orderID,
		ToStatus:    toStatus,
		SourceType:  "ADMIN",
		Description: stringPtrOrNil(reason),
		OccurredAt:  now,
	}

	notification := &entities.Notification{
		ID:        uuid.NewString(),
		UserID:    order.UserID,
		Channel:   "IN_APP",
		EventType: "order.status.changed",
		EventKey:  orderID,
		Title:     "Đơn hàng cập nhật trạng thái",
		Content:   fmt.Sprintf("Đơn hàng chuyển sang %s", toStatus),
		Status:    entities.NotificationStatusSent,
		IsRead:    false,
		SentAt:    &now,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.repo.UpdateOrderStatusTx(orderID, updates, tracking, notification)
	if err != nil {
		if strings.Contains(err.Error(), "INVALID_TRANSITION") {
			return nil, &apperrors.APIError{Status: http.StatusConflict, Code: "CONFLICT", Message: "invalid order status transition"}
		}
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to commit status update"}
	}

	return s.GetByIDAdmin(orderID)
}

func (s *Service) MarkPaymentResult(orderID string, success bool, reason string) (*entities.Order, *apperrors.APIError) {
	now := time.Now().UTC()

	order, err := s.repo.FindByID(orderID)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, apperrors.ErrNotFound
		}
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query order"}
	}

	updates := map[string]any{
		"updated_at": now,
	}
	var tracking *entities.OrderTrackingEvent

	notification := &entities.Notification{
		ID:        uuid.NewString(),
		UserID:    order.UserID,
		Channel:   "IN_APP",
		Status:    entities.NotificationStatusSent,
		IsRead:    false,
		SentAt:    &now,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if success {
		updates["payment_status"] = entities.PaymentStatusPaid
		if order.OrderStatus == entities.OrderStatusPendingPayment {
			updates["order_status"] = entities.OrderStatusPaid

			tracking = &entities.OrderTrackingEvent{
				ID:          uuid.NewString(),
				OrderID:     orderID,
				ToStatus:    entities.OrderStatusPaid,
				SourceType:  "PAYMENT_GATEWAY",
				Description: stringPtrOrNil("Payment succeeded"),
				OccurredAt:  now,
			}
		}

		notification.EventType = "payment.succeeded"
		notification.EventKey = orderID
		notification.Title = "Thanh toán thành công"
		notification.Content = "Thanh toán đơn hàng thành công"

	} else {
		updates["payment_status"] = entities.PaymentStatusFailed
		notification.EventType = "payment.failed"
		notification.EventKey = orderID
		notification.Title = "Thanh toán thất bại"
		notification.Content = reason
	}

	err = s.repo.MarkPaymentResultTx(orderID, updates, tracking, notification)
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to commit payment update"}
	}

	return s.GetByIDAdmin(orderID)
}

func (s *Service) getOrder(orderID string) (*entities.Order, error) {
	return s.repo.FindByID(orderID)
}

func stringPtrOrNil(v string) *string {
	trimmed := strings.TrimSpace(v)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func ParseOrderStatus(value string) (string, error) {
	up := strings.ToUpper(strings.TrimSpace(value))
	switch up {
	case entities.OrderStatusNew,
		entities.OrderStatusPendingPayment,
		entities.OrderStatusPaid,
		entities.OrderStatusProcessing,
		entities.OrderStatusShipping,
		entities.OrderStatusCompleted,
		entities.OrderStatusCancelled:
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

func cloneOrder(entity *entities.Order) *entities.Order {
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

func cloneOrderItem(entity *entities.OrderItem) *entities.OrderItem {
	if entity == nil {
		return nil
	}
	copy := *entity
	copy.CreatedAt = copy.CreatedAt.UTC()
	copy.UpdatedAt = copy.UpdatedAt.UTC()
	return &copy
}

func cloneTracking(entity *entities.OrderTrackingEvent) *entities.OrderTrackingEvent {
	if entity == nil {
		return nil
	}
	copy := *entity
	copy.OccurredAt = copy.OccurredAt.UTC()
	return &copy
}
