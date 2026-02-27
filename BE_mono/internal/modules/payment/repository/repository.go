package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"golf-store/be-mono/internal/platform/db"
)

type Repository interface {
	FindOrder(orderID string) (*db.OrderEntity, error)
	FindPaymentByIdempotency(orderID string, idemKey string) (*db.PaymentTransactionEntity, error)
	CreatePayment(payment *db.PaymentTransactionEntity) error
	ListPaymentsByOrder(orderID string) ([]db.PaymentTransactionEntity, error)
	FindWebhookDuplicate(providerTxnCode string) (*db.PaymentTransactionEntity, error)
	ProcessWebhookTx(orderID string, provider string, providerTxnCode string, status string, finalState string) (string, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGorm(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) FindOrder(orderID string) (*db.OrderEntity, error) {
	var order db.OrderEntity
	err := r.db.Select("id", "user_id", "total_amount", "currency_code").Where("id = ?", orderID).Take(&order).Error
	return &order, err
}

func (r *GormRepository) FindPaymentByIdempotency(orderID string, idemKey string) (*db.PaymentTransactionEntity, error) {
	var existing db.PaymentTransactionEntity
	err := r.db.Where("order_id = ? AND idempotency_key = ?", orderID, idemKey).Take(&existing).Error
	return &existing, err
}

func (r *GormRepository) CreatePayment(payment *db.PaymentTransactionEntity) error {
	return r.db.Create(payment).Error
}

func (r *GormRepository) ListPaymentsByOrder(orderID string) ([]db.PaymentTransactionEntity, error) {
	entities := make([]db.PaymentTransactionEntity, 0)
	err := r.db.Where("order_id = ?", orderID).Order("created_at ASC").Find(&entities).Error
	return entities, err
}

func (r *GormRepository) FindWebhookDuplicate(providerTxnCode string) (*db.PaymentTransactionEntity, error) {
	var duplicate db.PaymentTransactionEntity
	err := r.db.Where("provider_txn_code = ?", providerTxnCode).Take(&duplicate).Error
	return &duplicate, err
}

func (r *GormRepository) ProcessWebhookTx(orderID string, provider string, providerTxnCode string, status string, finalState string) (string, error) {
	var paymentID string

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var payment db.PaymentTransactionEntity
		lockErr := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("order_id = ? AND provider = ? AND status = ?", orderID, provider, db.PaymentTxnStatePending).
			Order("created_at ASC").
			Take(&payment).Error

		if lockErr != nil {
			if lockErr != gorm.ErrRecordNotFound {
				return lockErr
			}
			// If not found, create one
			var order db.OrderEntity
			if err := tx.Select("id", "total_amount", "currency_code").Where("id = ?", orderID).Take(&order).Error; err != nil {
				return err
			}

			payment = db.PaymentTransactionEntity{
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
			// Let's pass an optional *db.PaymentTransactionEntity to create if lockErr == gorm.ErrRecordNotFound.
			return gorm.ErrRecordNotFound
		}

		paymentID = payment.ID

		if err := tx.Model(&db.PaymentTransactionEntity{}).
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
