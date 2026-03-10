package service

import (
	"errors"
	"testing"
	"time"

	"gorm.io/gorm"

	"golf-store/be-mono/internal/modules/order/repository"
	"golf-store/be-mono/internal/platform/entities"
)

type fakeOrderRepo struct {
	order        *entities.Order
	expireCalled bool
}

func (f *fakeOrderRepo) FindIdempotentOrder(userID string, idemKey string) (*entities.Order, error) {
	return nil, gorm.ErrRecordNotFound
}

func (f *fakeOrderRepo) CreateOrderTx(order *entities.Order, items []entities.OrderItem, tracking *entities.OrderTrackingEvent, notification *entities.Notification, productUpdates map[string]int) error {
	return nil
}

func (f *fakeOrderRepo) List(filter repository.ListFilter) ([]entities.Order, int64, error) {
	return nil, 0, nil
}

func (f *fakeOrderRepo) FindByID(orderID string) (*entities.Order, error) {
	if f.order == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return f.order, nil
}

func (f *fakeOrderRepo) CancelOrderTx(orderID string, userID string, updates map[string]any, tracking *entities.OrderTrackingEvent, notification *entities.Notification) error {
	return nil
}

func (f *fakeOrderRepo) ExpirePendingPaymentTx(orderID string, updates map[string]any, tracking *entities.OrderTrackingEvent, notification *entities.Notification) error {
	if f.order == nil {
		return gorm.ErrRecordNotFound
	}
	f.expireCalled = true
	f.order.OrderStatus = updates["order_status"].(string)
	f.order.PaymentStatus = updates["payment_status"].(string)
	f.order.CancelReason = stringPtrOrNil(updates["cancel_reason"].(string))
	f.order.UpdatedAt = updates["updated_at"].(time.Time)
	return nil
}

func (f *fakeOrderRepo) UpdateOrderStatusTx(orderID string, updates map[string]any, tracking *entities.OrderTrackingEvent, notification *entities.Notification) error {
	return nil
}

func (f *fakeOrderRepo) MarkPaymentResultTx(orderID string, updates map[string]any, tracking *entities.OrderTrackingEvent, notification *entities.Notification) error {
	return nil
}

func (f *fakeOrderRepo) LoadTracking(orderID string) ([]entities.OrderTrackingEvent, error) {
	return nil, nil
}

func (f *fakeOrderRepo) FindUserByID(userID string) (*entities.User, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeOrderRepo) LockProduct(tx *gorm.DB, productID string) (*entities.Product, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeOrderRepo) GetDbConn() *gorm.DB {
	return nil
}

func TestExpireIfPastDueCancelsExpiredPendingOrder(t *testing.T) {
	pastDue := time.Now().UTC().Add(-1 * time.Minute)
	repo := &fakeOrderRepo{
		order: &entities.Order{
			ID:            "order-1",
			UserID:        "user-1",
			OrderStatus:   entities.OrderStatusPendingPayment,
			PaymentStatus: entities.PaymentStatusUnpaid,
			PaymentDueAt:  &pastDue,
		},
	}
	svc := New(repo)

	order, expired, apiErr := svc.ExpireIfPastDue("order-1")
	if apiErr != nil {
		t.Fatalf("expected no error, got %v", apiErr)
	}
	if !expired {
		t.Fatalf("expected order to expire")
	}
	if !repo.expireCalled {
		t.Fatalf("expected repository expiration to be called")
	}
	if order.OrderStatus != entities.OrderStatusCancelled {
		t.Fatalf("expected cancelled status, got %s", order.OrderStatus)
	}
	if order.PaymentStatus != entities.PaymentStatusFailed {
		t.Fatalf("expected failed payment status, got %s", order.PaymentStatus)
	}
}

func TestExpireIfPastDueSkipsFreshOrder(t *testing.T) {
	futureDue := time.Now().UTC().Add(10 * time.Minute)
	repo := &fakeOrderRepo{
		order: &entities.Order{
			ID:            "order-1",
			UserID:        "user-1",
			OrderStatus:   entities.OrderStatusPendingPayment,
			PaymentStatus: entities.PaymentStatusUnpaid,
			PaymentDueAt:  &futureDue,
		},
	}
	svc := New(repo)

	order, expired, apiErr := svc.ExpireIfPastDue("order-1")
	if apiErr != nil {
		t.Fatalf("expected no error, got %v", apiErr)
	}
	if expired {
		t.Fatalf("expected order not to expire")
	}
	if repo.expireCalled {
		t.Fatalf("did not expect repository expiration to be called")
	}
	if order.OrderStatus != entities.OrderStatusPendingPayment {
		t.Fatalf("expected pending status, got %s", order.OrderStatus)
	}
}
