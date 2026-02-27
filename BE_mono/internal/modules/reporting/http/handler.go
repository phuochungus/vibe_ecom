package http

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

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
	response.OK(c, gin.H{
		"from":                 out.From.Format(time.RFC3339Nano),
		"to":                   out.To.Format(time.RFC3339Nano),
		"gross_revenue":        utils.ToAmountString(out.GrossRevenue),
		"refund_amount":        utils.ToAmountString(out.RefundAmount),
		"net_revenue":          utils.ToAmountString(out.NetRevenue),
		"completed_orders":     out.CompletedOrders,
		"payment_success_rate": out.PaymentSuccessRate,
	})
}

func (h *Handler) RevenueOrders(c *gin.Context) {
	from, to := parseRange(c)
	page := parseIntDefault(c.Query("page"), 1)
	pageSize := parseIntDefault(c.Query("page_size"), 20)

	out := h.reporting.Orders(from, to, page, pageSize)
	items := make([]gin.H, 0, len(out.Items))
	for _, order := range out.Items {
		items = append(items, gin.H{
			"id":              order.ID,
			"order_code":      order.OrderCode,
			"order_status":    order.OrderStatus,
			"payment_status":  order.PaymentStatus,
			"subtotal_amount": utils.ToAmountString(order.SubtotalAmount),
			"discount_amount": utils.ToAmountString(order.DiscountAmount),
			"shipping_fee":    utils.ToAmountString(order.ShippingFee),
			"total_amount":    utils.ToAmountString(order.TotalAmount),
			"placed_at":       order.PlacedAt.Format(time.RFC3339Nano),
		})
	}

	response.OK(c, gin.H{
		"items": items,
		"pagination": gin.H{
			"page":        out.Page,
			"page_size":   out.PageSize,
			"total":       out.Total,
			"total_pages": out.TotalPages,
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
