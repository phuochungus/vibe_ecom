package http

import (
	"fmt"
	nethttp "net/http"

	"github.com/gin-gonic/gin"

	"golf-store/be/internal/platform/observability"
	gatewaymsg "golf-store/be/internal/services/gateway/messaging"
)

type Handler struct {
	serviceName string
	publisher   gatewaymsg.CommandPublisher
}

func NewHandler(serviceName string, publisher gatewaymsg.CommandPublisher) *Handler {
	return &Handler{
		serviceName: serviceName,
		publisher:   publisher,
	}
}

type createOrderRequest struct {
	UserID string                 `json:"userId"`
	Items  []createOrderLineInput `json:"items"`
}

type createOrderLineInput struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(nethttp.StatusOK, gin.H{
		"service": h.serviceName,
		"status":  "ok",
	})
}

func (h *Handler) Ready(c *gin.Context) {
	c.JSON(nethttp.StatusOK, gin.H{
		"service": h.serviceName,
		"ready":   true,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var payload map[string]any
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	h.acceptCommand(c, "auth.login.requested", payload)
}

func (h *Handler) ListProducts(c *gin.Context) {
	c.JSON(nethttp.StatusOK, gin.H{
		"items": []any{},
		"meta": gin.H{
			"page":  1,
			"limit": 20,
			"total": 0,
		},
	})
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if len(req.Items) == 0 {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": "order must include at least one item"})
		return
	}

	h.acceptCommand(c, "order.create.requested", req)
}

func (h *Handler) GetOrderTracking(c *gin.Context) {
	c.JSON(nethttp.StatusOK, gin.H{
		"orderCode": c.Param("orderCode"),
		"timeline":  []any{},
	})
}

func (h *Handler) ListNotifications(c *gin.Context) {
	c.JSON(nethttp.StatusOK, gin.H{
		"items": []any{},
	})
}

func (h *Handler) ReceivePaymentWebhook(c *gin.Context) {
	provider := c.Param("provider")

	var payload map[string]any
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	commandType := fmt.Sprintf("payment.webhook.%s.received", provider)
	h.acceptCommand(c, commandType, payload)
}

func (h *Handler) AdminListOrders(c *gin.Context) {
	c.JSON(nethttp.StatusOK, gin.H{
		"items": []any{},
	})
}

func (h *Handler) AdminRevenueSummary(c *gin.Context) {
	c.JSON(nethttp.StatusOK, gin.H{
		"currency": "VND",
		"gross":    0,
		"refund":   0,
		"net":      0,
	})
}

func (h *Handler) acceptCommand(c *gin.Context, commandType string, payload any) {
	correlationID := observability.CorrelationIDFromContext(c)

	envelope, err := h.publisher.PublishCommand(c.Request.Context(), correlationID, commandType, payload)
	if err != nil {
		c.JSON(nethttp.StatusInternalServerError, gin.H{"error": "failed to enqueue command"})
		return
	}

	c.JSON(nethttp.StatusAccepted, gin.H{
		"status":        "accepted",
		"messageType":   envelope.Type,
		"commandId":     envelope.ID,
		"correlationId": envelope.CorrelationID,
	})
}
