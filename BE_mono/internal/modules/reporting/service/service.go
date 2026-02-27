package service

import (
	"fmt"
	"time"

	"golf-store/be-mono/internal/modules/reporting/repository"
	"golf-store/be-mono/internal/platform/db"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
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
	grossRevenue, _ := s.repo.SumGrossRevenue(from, to)
	completedOrders, _ := s.repo.CountCompletedOrders(from, to)
	refundAmount, _ := s.repo.SumRefundAmount(from, to)
	paymentTotal, _ := s.repo.CountPaymentTotal(from, to)
	paymentSuccess, _ := s.repo.CountPaymentSuccess(from, to)

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

	entities, total, err := s.repo.ListCompletedOrders(from, to, page, pageSize)
	if err != nil {
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
