package http

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"golf-store/be-mono/internal/modules/reporting/dto"
	reportsvc "golf-store/be-mono/internal/modules/reporting/service"
	"golf-store/be-mono/internal/shared/response"
	"golf-store/be-mono/internal/shared/utils"
)

type Handler struct {
	reporting *reportsvc.Service
}

func New(reporting *reportsvc.Service) *Handler {
	return &Handler{reporting: reporting}
}

func (h *Handler) RegisterAdmin(rg *gin.RouterGroup) {
	rg.GET("/revenue/summary", h.Summary)
	rg.GET("/revenue/orders", h.RevenueOrders)
}

func (h *Handler) Summary(c *gin.Context) {
	from, to := parseRange(c)
	out := h.reporting.Summary(from, to)
	response.OK(c, dto.SummaryResponseDTO{
		From:               out.From.Format(time.RFC3339Nano),
		To:                 out.To.Format(time.RFC3339Nano),
		GrossRevenue:       utils.ToAmountString(out.GrossRevenue),
		RefundAmount:       utils.ToAmountString(out.RefundAmount),
		NetRevenue:         utils.ToAmountString(out.NetRevenue),
		CompletedOrders:    out.CompletedOrders,
		PaymentSuccessRate: out.PaymentSuccessRate,
	})
}

func (h *Handler) RevenueOrders(c *gin.Context) {
	from, to := parseRange(c)
	page := parseIntDefault(c.Query("page"), 1)
	pageSize := parseIntDefault(c.Query("page_size"), 20)

	out := h.reporting.Orders(from, to, page, pageSize)
	items := make([]dto.RevenueOrderItemDTO, 0, len(out.Items))
	for _, order := range out.Items {
		items = append(items, dto.RevenueOrderItemDTO{
			ID:             order.ID,
			OrderCode:      order.OrderCode,
			OrderStatus:    order.OrderStatus,
			PaymentStatus:  order.PaymentStatus,
			SubtotalAmount: utils.ToAmountString(order.SubtotalAmount),
			DiscountAmount: utils.ToAmountString(order.DiscountAmount),
			ShippingFee:    utils.ToAmountString(order.ShippingFee),
			TotalAmount:    utils.ToAmountString(order.TotalAmount),
			PlacedAt:       order.PlacedAt.Format(time.RFC3339Nano),
		})
	}

	response.OK(c, dto.RevenueOrdersResponseDTO{
		Items: items,
		Pagination: dto.PaginationDTO{
			Page:       out.Page,
			PageSize:   out.PageSize,
			Total:      out.Total,
			TotalPages: out.TotalPages,
		},
	})
}

func parseRange(c *gin.Context) (time.Time, time.Time) {
	now := time.Now().UTC()
	defaultFrom := now.AddDate(0, -1, 0)
	from := parseTimeOrDefault(c.Query("from"), defaultFrom)
	to := parseTimeOrDefault(c.Query("to"), now)
	if !from.Before(to) {
		to = from.Add(24 * time.Hour)
	}
	return from, to
}

func parseTimeOrDefault(value string, fallback time.Time) time.Time {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return fallback
	}
	return t.UTC()
}

func parseIntDefault(value string, fallback int) int {
	v, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || v <= 0 {
		return fallback
	}
	return v
}
