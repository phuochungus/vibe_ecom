package service

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"golf-store/be-mono/internal/platform/db"
)

type Service struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Service {
	return &Service{db: db}
}

type SummaryOutput struct {
	From               time.Time
	To                 time.Time
	GrossRevenue       int64
	RefundAmount       int64
	NetRevenue         int64
	CompletedOrders    int
	PaymentSuccessRate string
}

type OrdersOutput struct {
	Items      []*db.OrderEntity
	Page       int
	PageSize   int
	Total      int
	TotalPages int
}

func (s *Service) Summary(from time.Time, to time.Time) SummaryOutput {
	var grossRevenue int64
	_ = s.db.Model(&db.OrderEntity{}).
		Where("order_status = ? AND placed_at >= ? AND placed_at < ?", db.OrderStatusCompleted, from, to).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&grossRevenue).Error

	var completedOrders int64
	_ = s.db.Model(&db.OrderEntity{}).
		Where("order_status = ? AND placed_at >= ? AND placed_at < ?", db.OrderStatusCompleted, from, to).
		Count(&completedOrders).Error

	var refundAmount int64
	_ = s.db.Model(&db.PaymentTransactionEntity{}).
		Where("txn_type = ? AND status = ? AND created_at >= ? AND created_at < ?", db.PaymentTxnTypeRefund, db.PaymentTxnStateSuccess, from, to).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&refundAmount).Error

	var paymentTotal int64
	_ = s.db.Model(&db.PaymentTransactionEntity{}).
		Where("txn_type = ? AND created_at >= ? AND created_at < ?", db.PaymentTxnTypePayment, from, to).
		Count(&paymentTotal).Error

	var paymentSuccess int64
	_ = s.db.Model(&db.PaymentTransactionEntity{}).
		Where("txn_type = ? AND status = ? AND created_at >= ? AND created_at < ?", db.PaymentTxnTypePayment, db.PaymentTxnStateSuccess, from, to).
		Count(&paymentSuccess).Error

	rate := "0.00"
	if paymentTotal > 0 {
		rate = fmt.Sprintf("%.2f", float64(paymentSuccess)/float64(paymentTotal))
	}

	return SummaryOutput{
		From:               from,
		To:                 to,
		GrossRevenue:       grossRevenue,
		RefundAmount:       refundAmount,
		NetRevenue:         grossRevenue - refundAmount,
		CompletedOrders:    int(completedOrders),
		PaymentSuccessRate: rate,
	}
}

func (s *Service) Orders(from time.Time, to time.Time, page int, pageSize int) OrdersOutput {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	query := s.db.Model(&db.OrderEntity{}).
		Where("order_status = ? AND placed_at >= ? AND placed_at < ?", db.OrderStatusCompleted, from, to)

	var total int64
	_ = query.Count(&total).Error

	offset := (page - 1) * pageSize
	entities := make([]db.OrderEntity, 0)
	if err := query.Order("placed_at DESC").Limit(pageSize).Offset(offset).Find(&entities).Error; err != nil {
		return OrdersOutput{Items: []*db.OrderEntity{}, Page: page, PageSize: pageSize, Total: int(total), TotalPages: calcTotalPages(int(total), pageSize)}
	}

	items := make([]*db.OrderEntity, 0, len(entities))
	for i := range entities {
		items = append(items, cloneOrder(&entities[i]))
	}

	return OrdersOutput{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: calcTotalPages(int(total), pageSize),
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
