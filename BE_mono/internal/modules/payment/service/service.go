package service

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	ordersvc "golf-store/be-mono/internal/modules/order/service"
	"golf-store/be-mono/internal/modules/payment/repository"
	"golf-store/be-mono/internal/platform/db"
	apperrors "golf-store/be-mono/internal/shared/errors"
)

type Service struct {
	repo   repository.Repository
	orders *ordersvc.Service
}

func New(repo repository.Repository, orders *ordersvc.Service) *Service {
	return &Service{repo: repo, orders: orders}
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

	order, err := s.repo.FindOrder(input.OrderID)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, apperrors.ErrNotFound
		}
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query order"}
	}
	if order.UserID != input.UserID {
		return nil, apperrors.ErrNotFound
	}

	existing, err := s.repo.FindPaymentByIdempotency(input.OrderID, idemKey)
	if err == nil && existing != nil && existing.ID != "" {
		return &CreatePaymentOutput{
			PaymentID:   existing.ID,
			Status:      existing.Status,
			CheckoutURL: fakeCheckoutURL(input.Provider, existing.ID),
		}, nil
	}
	if err != nil && !strings.Contains(err.Error(), "record not found") {
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
	if err := s.repo.CreatePayment(&payment); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusConflict, Code: "PAYMENT_DUPLICATE", Message: "Duplicate payment request"}
	}

	return &CreatePaymentOutput{
		PaymentID:   paymentID,
		Status:      db.PaymentTxnStatePending,
		CheckoutURL: fakeCheckoutURL(provider, paymentID),
	}, nil
}

func (s *Service) ListByOrderForUser(orderID string, userID string) ([]*db.PaymentTransactionEntity, *apperrors.APIError) {
	order, err := s.repo.FindOrder(orderID)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, apperrors.ErrNotFound
		}
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query order"}
	}
	if order.UserID != userID {
		return nil, apperrors.ErrNotFound
	}

	entities, err := s.repo.ListPaymentsByOrder(orderID)
	if err != nil {
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

	duplicate, err := s.repo.FindWebhookDuplicate(providerTxnCode)
	if err == nil && duplicate != nil && duplicate.ID != "" {
		deduplicated = true
		paymentID = duplicate.ID
	} else if err != nil && !strings.Contains(err.Error(), "record not found") {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query webhook dedupe"}
	}

	if !deduplicated {
		if status == "SUCCESS" || status == "PAID" {
			finalState = db.PaymentTxnStateSuccess
		} else {
			finalState = db.PaymentTxnStateFailed
		}

		paymentID, err = s.repo.ProcessWebhookTx(orderID, provider, providerTxnCode, status, finalState)
		if err != nil {
			if strings.Contains(err.Error(), "record not found") {
				return nil, apperrors.ErrNotFound
			}
			return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to process payment webhook"}
		}
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
