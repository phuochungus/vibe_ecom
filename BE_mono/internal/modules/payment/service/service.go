package service

import (
	"database/sql"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"

	ordersvc "golf-store/be-mono/internal/modules/order/service"
	apperrors "golf-store/be-mono/internal/shared/errors"
	"golf-store/be-mono/internal/shared/model"
)

type Service struct {
	db     *sql.DB
	orders *ordersvc.Service
}

func New(db *sql.DB, orders *ordersvc.Service) *Service {
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
	Status      model.PaymentTxnState
	CheckoutURL string
}

func (s *Service) Create(input CreatePaymentInput) (*CreatePaymentOutput, *apperrors.APIError) {
	if strings.TrimSpace(input.IdempotencyKey) == "" {
		return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: "Idempotency-Key is required"}
	}

	var orderUserID string
	var amount int64
	var currency string
	err := s.db.QueryRow(
		`SELECT user_id, total_amount, currency_code FROM orders WHERE id = ? LIMIT 1`,
		input.OrderID,
	).Scan(&orderUserID, &amount, &currency)
	if err == sql.ErrNoRows {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query order"}
	}
	if orderUserID != input.UserID {
		return nil, apperrors.ErrNotFound
	}

	var existingID string
	var existingStatus string
	err = s.db.QueryRow(
		`SELECT id, status
		   FROM payment_transactions
		  WHERE order_id = ? AND idempotency_key = ?
		  LIMIT 1`,
		input.OrderID, strings.TrimSpace(input.IdempotencyKey),
	).Scan(&existingID, &existingStatus)
	if err == nil {
		return &CreatePaymentOutput{
			PaymentID:   existingID,
			Status:      model.PaymentTxnState(existingStatus),
			CheckoutURL: fakeCheckoutURL(input.Provider, existingID),
		}, nil
	}
	if err != nil && err != sql.ErrNoRows {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query payment idempotency"}
	}

	now := time.Now().UTC()
	paymentID := uuid.NewString()
	provider := strings.ToUpper(strings.TrimSpace(input.Provider))
	if provider == "" {
		provider = "PAYOS"
	}

	if _, err := s.db.Exec(
		`INSERT INTO payment_transactions
		  (id, order_id, txn_type, provider, provider_txn_code, idempotency_key, amount, currency_code, status, provider_response, created_at, updated_at)
		 VALUES (?, ?, 'PAYMENT', ?, NULL, ?, ?, ?, 'PENDING', NULL, ?, ?)`,
		paymentID, input.OrderID, provider, strings.TrimSpace(input.IdempotencyKey), amount, currency, now, now,
	); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusConflict, Code: "PAYMENT_DUPLICATE", Message: "Duplicate payment request"}
	}

	return &CreatePaymentOutput{
		PaymentID:   paymentID,
		Status:      model.PaymentTxnStatePending,
		CheckoutURL: fakeCheckoutURL(provider, paymentID),
	}, nil
}

func (s *Service) ListByOrderForUser(orderID string, userID string) ([]*model.PaymentTransaction, *apperrors.APIError) {
	var orderUserID string
	if err := s.db.QueryRow(`SELECT user_id FROM orders WHERE id = ?`, orderID).Scan(&orderUserID); err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query order"}
	}
	if orderUserID != userID {
		return nil, apperrors.ErrNotFound
	}

	rows, err := s.db.Query(
		`SELECT id, order_id, txn_type, provider, provider_txn_code, idempotency_key, amount, currency_code, status, provider_response, created_at, updated_at
		   FROM payment_transactions
		  WHERE order_id = ?
		  ORDER BY created_at ASC`,
		orderID,
	)
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query payments"}
	}
	defer rows.Close()

	list := make([]*model.PaymentTransaction, 0)
	for rows.Next() {
		p, err := scanPayment(rows)
		if err != nil {
			continue
		}
		list = append(list, p)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt.Before(list[j].CreatedAt)
	})
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

	tx, err := s.db.Begin()
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to begin transaction"}
	}
	defer tx.Rollback()

	var duplicateID string
	if err := tx.QueryRow(
		`SELECT id FROM payment_transactions WHERE provider_txn_code = ? LIMIT 1`,
		providerTxnCode,
	).Scan(&duplicateID); err == nil {
		if err := tx.Commit(); err != nil {
			return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to commit webhook dedupe"}
		}
		return map[string]any{"status": "deduplicated", "payment_id": duplicateID, "provider_txn_code": providerTxnCode}, nil
	}

	paymentID := ""
	var amount int64
	var currency string
	row := tx.QueryRow(
		`SELECT id, amount, currency_code
		   FROM payment_transactions
		  WHERE order_id = ? AND provider = ? AND status = 'PENDING'
		  ORDER BY created_at ASC
		  LIMIT 1
		  FOR UPDATE`,
		orderID, provider,
	)
	if err := row.Scan(&paymentID, &amount, &currency); err == sql.ErrNoRows {
		paymentID = uuid.NewString()
		if err := tx.QueryRow(`SELECT total_amount, currency_code FROM orders WHERE id = ?`, orderID).Scan(&amount, &currency); err != nil {
			if err == sql.ErrNoRows {
				return nil, apperrors.ErrNotFound
			}
			return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query order amount"}
		}
		if _, err := tx.Exec(
			`INSERT INTO payment_transactions
			  (id, order_id, txn_type, provider, provider_txn_code, idempotency_key, amount, currency_code, status, provider_response, created_at, updated_at)
			 VALUES (?, ?, 'PAYMENT', ?, ?, NULL, ?, ?, 'PENDING', NULL, UTC_TIMESTAMP(3), UTC_TIMESTAMP(3))`,
			paymentID, orderID, provider, providerTxnCode, amount, currency,
		); err != nil {
			return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to create payment transaction"}
		}
	} else if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query payment transaction"}
	}

	finalState := model.PaymentTxnStateFailed
	if status == "SUCCESS" || status == "PAID" {
		finalState = model.PaymentTxnStateSuccess
	}

	if _, err := tx.Exec(
		`UPDATE payment_transactions
		    SET provider_txn_code = ?, provider_response = ?, status = ?, updated_at = UTC_TIMESTAMP(3)
		  WHERE id = ?`,
		providerTxnCode, status, finalState, paymentID,
	); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to update payment transaction"}
	}

	if err := tx.Commit(); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to commit payment webhook"}
	}

	_, orderErr := s.orders.MarkPaymentResult(orderID, finalState == model.PaymentTxnStateSuccess, "Payment webhook processed")
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

func scanPayment(scanner interface {
	Scan(dest ...any) error
}) (*model.PaymentTransaction, error) {
	payment := &model.PaymentTransaction{}
	var txnType, status string
	var providerTxnCode sql.NullString
	var idempotency sql.NullString
	var providerResp sql.NullString

	if err := scanner.Scan(
		&payment.ID,
		&payment.OrderID,
		&txnType,
		&payment.Provider,
		&providerTxnCode,
		&idempotency,
		&payment.Amount,
		&payment.CurrencyCode,
		&status,
		&providerResp,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	); err != nil {
		return nil, err
	}
	payment.TxnType = model.PaymentTxnType(txnType)
	payment.Status = model.PaymentTxnState(status)
	if providerTxnCode.Valid {
		payment.ProviderTxnCode = providerTxnCode.String
	}
	if idempotency.Valid {
		payment.IdempotencyKey = idempotency.String
	}
	if providerResp.Valid {
		payment.ProviderResponse = providerResp.String
	}
	return payment, nil
}

func fakeCheckoutURL(provider string, paymentID string) string {
	if provider == "" {
		provider = "PAYOS"
	}
	return fmt.Sprintf("https://checkout.local/%s/%s", strings.ToLower(provider), paymentID)
}
