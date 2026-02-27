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

	ordersvc "golf-store/be-mono/internal/modules/order/service"
	"golf-store/be-mono/internal/platform/db"
	apperrors "golf-store/be-mono/internal/shared/errors"
)

type Service struct {
	db     *gorm.DB
	orders *ordersvc.Service
}

func New(db *gorm.DB, orders *ordersvc.Service) *Service {
	return &Service{db: db, orders: orders}
}

type CreatePaymentInput struct {
	UserID         string
	OrderID        string
	IdempotencyKey string
	Provider       string
	ReturnURL      string
	CancelURL      string
}

type CreatePaymentOutput struct {
	PaymentID   string
	Status      string
	CheckoutURL string
}

func (s *Service) Create(input CreatePaymentInput) (*CreatePaymentOutput, *apperrors.APIError) {
	idemKey := strings.TrimSpace(input.IdempotencyKey)
	if idemKey == "" {
		return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: "Idempotency-Key is required"}
	}

	var order db.OrderEntity
	err := s.db.Select("id", "user_id", "total_amount", "currency_code").Where("id = ?", input.OrderID).Take(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query order"}
	}
	if order.UserID != input.UserID {
		return nil, apperrors.ErrNotFound
	}

	var existing db.PaymentTransactionEntity
	err = s.db.Where("order_id = ? AND idempotency_key = ?", input.OrderID, idemKey).Take(&existing).Error
	if err == nil {
		return &CreatePaymentOutput{
			PaymentID:   existing.ID,
			Status:      existing.Status,
			CheckoutURL: fakeCheckoutURL(input.Provider, existing.ID),
		}, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query payment idempotency"}
	}

	now := time.Now().UTC()
	paymentID := uuid.NewString()
	provider := strings.ToUpper(strings.TrimSpace(input.Provider))
	if provider == "" {
		provider = "PAYOS"
	}

	payment := db.PaymentTransactionEntity{
		ID:             paymentID,
		OrderID:        input.OrderID,
		TxnType:        db.PaymentTxnTypePayment,
		Provider:       provider,
		IdempotencyKey: &idemKey,
		Amount:         order.TotalAmount,
		CurrencyCode:   order.CurrencyCode,
		Status:         db.PaymentTxnStatePending,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if err := s.db.Create(&payment).Error; err != nil {
		return nil, &apperrors.APIError{Status: http.StatusConflict, Code: "PAYMENT_DUPLICATE", Message: "Duplicate payment request"}
	}

	return &CreatePaymentOutput{
		PaymentID:   paymentID,
		Status:      db.PaymentTxnStatePending,
		CheckoutURL: fakeCheckoutURL(provider, paymentID),
	}, nil
}

func (s *Service) ListByOrderForUser(orderID string, userID string) ([]*db.PaymentTransactionEntity, *apperrors.APIError) {
	var order db.OrderEntity
	if err := s.db.Select("id", "user_id").Where("id = ?", orderID).Take(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query order"}
	}
	if order.UserID != userID {
		return nil, apperrors.ErrNotFound
	}

	entities := make([]db.PaymentTransactionEntity, 0)
	if err := s.db.Where("order_id = ?", orderID).Order("created_at ASC").Find(&entities).Error; err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query payments"}
	}

	list := make([]*db.PaymentTransactionEntity, 0, len(entities))
	for i := range entities {
		list = append(list, clonePayment(&entities[i]))
	}
	return list, nil
}

func (s *Service) ProcessWebhook(provider string, payload map[string]any) (map[string]any, *apperrors.APIError) {
	provider = strings.ToUpper(strings.TrimSpace(provider))
	orderID, _ := payload["order_id"].(string)
	if strings.TrimSpace(orderID) == "" {
		return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: "order_id is required"}
	}

	statusRaw, _ := payload["status"].(string)
	status := strings.ToUpper(strings.TrimSpace(statusRaw))
	if status == "" {
		status = "SUCCESS"
	}

	providerTxnCode, _ := payload["provider_txn_code"].(string)
	if strings.TrimSpace(providerTxnCode) == "" {
		providerTxnCode = uuid.NewString()
	}
	providerTxnCode = strings.TrimSpace(providerTxnCode)

	paymentID := ""
	finalState := db.PaymentTxnStateFailed
	deduplicated := false

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var duplicate db.PaymentTransactionEntity
		if err := tx.Where("provider_txn_code = ?", providerTxnCode).Take(&duplicate).Error; err == nil {
			deduplicated = true
			paymentID = duplicate.ID
			return nil
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query webhook dedupe"}
		}

		var payment db.PaymentTransactionEntity
		lockErr := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("order_id = ? AND provider = ? AND status = ?", orderID, provider, db.PaymentTxnStatePending).
			Order("created_at ASC").
			Take(&payment).Error
		if lockErr != nil {
			if !errors.Is(lockErr, gorm.ErrRecordNotFound) {
				return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query payment transaction"}
			}

			var order db.OrderEntity
			if err := tx.Select("id", "total_amount", "currency_code").Where("id = ?", orderID).Take(&order).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return apperrors.ErrNotFound
				}
				return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query order amount"}
			}

			payment = db.PaymentTransactionEntity{
				ID:              uuid.NewString(),
				OrderID:         orderID,
				TxnType:         db.PaymentTxnTypePayment,
				Provider:        provider,
				ProviderTxnCode: &providerTxnCode,
				Amount:          order.TotalAmount,
				CurrencyCode:    order.CurrencyCode,
				Status:          db.PaymentTxnStatePending,
				CreatedAt:       time.Now().UTC(),
				UpdatedAt:       time.Now().UTC(),
			}
			if err := tx.Create(&payment).Error; err != nil {
				return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to create payment transaction"}
			}
		}

		paymentID = payment.ID
		if status == "SUCCESS" || status == "PAID" {
			finalState = db.PaymentTxnStateSuccess
		} else {
			finalState = db.PaymentTxnStateFailed
		}

		if err := tx.Model(&db.PaymentTransactionEntity{}).
			Where("id = ?", paymentID).
			Updates(map[string]any{
				"provider_txn_code": providerTxnCode,
				"provider_response": status,
				"status":            finalState,
				"updated_at":        time.Now().UTC(),
			}).Error; err != nil {
			return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to update payment transaction"}
		}

		return nil
	})
	if err != nil {
		if apiErr, ok := err.(*apperrors.APIError); ok {
			return nil, apiErr
		}
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to commit payment webhook"}
	}

	if deduplicated {
		return map[string]any{"status": "deduplicated", "payment_id": paymentID, "provider_txn_code": providerTxnCode}, nil
	}

	_, orderErr := s.orders.MarkPaymentResult(orderID, finalState == db.PaymentTxnStateSuccess, "Payment webhook processed")
	if orderErr != nil {
		return nil, orderErr
	}

	return map[string]any{
		"status":            "processed",
		"order_id":          orderID,
		"payment_id":        paymentID,
		"provider_txn_code": providerTxnCode,
		"payment_state":     finalState,
	}, nil
}

func clonePayment(entity *db.PaymentTransactionEntity) *db.PaymentTransactionEntity {
	if entity == nil {
		return nil
	}
	copy := *entity
	copy.CreatedAt = copy.CreatedAt.UTC()
	copy.UpdatedAt = copy.UpdatedAt.UTC()
	return &copy
}

func fakeCheckoutURL(provider string, paymentID string) string {
	if provider == "" {
		provider = "PAYOS"
	}
	return fmt.Sprintf("https://checkout.local/%s/%s", strings.ToLower(provider), paymentID)
}
