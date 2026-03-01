package repository

import (
	"time"

	"gorm.io/gorm"

	entities "golf-store/be-mono/internal/platform/entities"
)

type Repository interface {
	SumGrossRevenue(from, to time.Time) (int64, error)
	CountCompletedOrders(from, to time.Time) (int64, error)
	SumRefundAmount(from, to time.Time) (int64, error)
	CountPaymentTotal(from, to time.Time) (int64, error)
	CountPaymentSuccess(from, to time.Time) (int64, error)
	ListCompletedOrders(from, to time.Time, page, pageSize int) ([]*entities.Order, int64, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGorm(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) SumGrossRevenue(from, to time.Time) (int64, error) {
	var result int64
	err := r.db.Model(&entities.Order{}).
		Where("order_status = ? AND placed_at >= ? AND placed_at < ?", entities.OrderStatusCompleted, from, to).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&result).Error
	return result, err
}

func (r *GormRepository) CountCompletedOrders(from, to time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&entities.Order{}).
		Where("order_status = ? AND placed_at >= ? AND placed_at < ?", entities.OrderStatusCompleted, from, to).
		Count(&count).Error
	return count, err
}

func (r *GormRepository) SumRefundAmount(from, to time.Time) (int64, error) {
	var result int64
	err := r.db.Model(&entities.PaymentTransaction{}).
		Where("txn_type = ? AND status = ? AND created_at >= ? AND created_at < ?", entities.PaymentTxnTypeRefund, entities.PaymentTxnStateSuccess, from, to).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&result).Error
	return result, err
}

func (r *GormRepository) CountPaymentTotal(from, to time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&entities.PaymentTransaction{}).
		Where("txn_type = ? AND created_at >= ? AND created_at < ?", entities.PaymentTxnTypePayment, from, to).
		Count(&count).Error
	return count, err
}

func (r *GormRepository) CountPaymentSuccess(from, to time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&entities.PaymentTransaction{}).
		Where("txn_type = ? AND status = ? AND created_at >= ? AND created_at < ?", entities.PaymentTxnTypePayment, entities.PaymentTxnStateSuccess, from, to).
		Count(&count).Error
	return count, err
}

func (r *GormRepository) ListCompletedOrders(from, to time.Time, page, pageSize int) ([]*entities.Order, int64, error) {
	query := r.db.Model(&entities.Order{}).
		Where("order_status = ? AND placed_at >= ? AND placed_at < ?", entities.OrderStatusCompleted, from, to)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	entities := make([]*entities.Order, 0)
	if err := query.Order("placed_at DESC").Limit(pageSize).Offset(offset).Find(&entities).Error; err != nil {
		return nil, total, err
	}

	return entities, total, nil
}
