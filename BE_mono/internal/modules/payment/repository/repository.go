package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"golf-store/be-mono/internal/platform/db"
	entities "golf-store/be-mono/internal/platform/entities"
)

type Repository interface {
	FindOrder(orderID string) (*entities.Order, error)
	FindPaymentByIdempotency(orderID string, idemKey string) (*entities.PaymentTransaction, error)
	CreatePayment(payment *entities.PaymentTransaction) error
	ListPaymentsByOrder(orderID string) ([]entities.PaymentTransaction, error)
	FindWebhookDuplicate(providerTxnCode string) (*entities.PaymentTransaction, error)
	ProcessWebhookTx(orderID string, provider string, providerTxnCode string, status string, finalState string) (string, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGorm(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) FindOrder(orderID string) (*entities.Order, error) {
	var order entities.Order
	err := r.db.Select("id", "user_id", "total_amount", "currency_code").Where("id = ?", orderID).Take(&order).Error
	return &order, err
}

func (r *GormRepository) FindPaymentByIdempotency(orderID string, idemKey string) (*entities.PaymentTransaction, error) {
	var existing entities.PaymentTransaction
	err := r.db.Where("order_id = ? AND idempotency_key = ?", orderID, idemKey).Take(&existing).Error
	return &existing, err
}

func (r *GormRepository) CreatePayment(payment *entities.PaymentTransaction) error {
	return r.db.Create(payment).Error
}

func (r *GormRepository) ListPaymentsByOrder(orderID string) ([]entities.PaymentTransaction, error) {
	rows := make([]entities.PaymentTransaction, 0)
	err := r.db.Where("order_id = ?", orderID).Order("created_at ASC").Find(&rows).Error
	return rows, err
}

func (r *GormRepository) FindWebhookDuplicate(providerTxnCode string) (*entities.PaymentTransaction, error) {
	var duplicate entities.PaymentTransaction
	err := r.db.Where("provider_txn_code = ?", providerTxnCode).Take(&duplicate).Error
	return &duplicate, err
}

func (r *GormRepository) ProcessWebhookTx(orderID string, provider string, providerTxnCode string, status string, finalState string) (string, error) {
	var paymentID string

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var payment entities.PaymentTransaction
		lockErr := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("order_id = ? AND provider = ? AND status = ?", orderID, provider, db.PaymentTxnStatePending).
			Order("created_at ASC").
			Take(&payment).Error

		if lockErr != nil {
			if lockErr != gorm.ErrRecordNotFound {
				return lockErr
			}
			// If not found, create one
			var order entities.Order
			if err := tx.Select("id", "total_amount", "currency_code").Where("id = ?", orderID).Take(&order).Error; err != nil {
				return err
			}

			payment = entities.PaymentTransaction{
				ID:              providerTxnCode, // Assuming providerTxnCode is used as ID or injected later, handled by service if not uuid
				OrderID:         orderID,
				TxnType:         db.PaymentTxnTypePayment,
				Provider:        provider,
				ProviderTxnCode: &providerTxnCode,
				Amount:          order.TotalAmount,
				CurrencyCode:    order.CurrencyCode,
				Status:          db.PaymentTxnStatePending,
			}
			// The ID generation for new payment in webhook is handled in the service via uuid.NewString(), we will need to inject it or refactor. For now, returning the logic to service might be cleaner if ID generation is there, OR pass payment struct in from service.
			// Let's adjust this to take a pre-built payment entity to insert if needed, or return special error to let service handle it.

			// Actually, let's look at the original service code:
			// It generates a UUID on the fly and creates the payment in the transaction.
			// Let's pass an optional *entities.PaymentTransaction to create if lockErr == gorm.ErrRecordNotFound.
			return gorm.ErrRecordNotFound
		}

		paymentID = payment.ID

		if err := tx.Model(&entities.PaymentTransaction{}).
			Where("id = ?", paymentID).
			Updates(map[string]any{
				"provider_txn_code": providerTxnCode,
				"provider_response": status,
				"status":            finalState,
			}).Error; err != nil {
			return err
		}

		return nil
	})

	return paymentID, err
}
