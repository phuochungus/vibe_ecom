package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	paysvc "golf-store/be-mono/internal/modules/payment/service"
	"golf-store/be-mono/internal/shared/middleware"
	"golf-store/be-mono/internal/shared/response"
	"golf-store/be-mono/internal/shared/utils"
)

type Handler struct {
	payments *paysvc.Service
}

func New(payments *paysvc.Service) *Handler {
	return &Handler{payments: payments}
}

func (h *Handler) RegisterUser(rg *gin.RouterGroup) {
	rg.POST("/orders/:order_id/payments", h.CreatePayment)
	rg.GET("/orders/:order_id/payments", h.ListPayments)
}

func (h *Handler) RegisterWebhook(rg *gin.RouterGroup) {
	rg.POST("/webhooks/payments/:provider", h.ReceivePaymentWebhook)
}

type createPaymentRequest struct {
	Provider  string `json:"provider"`
	ReturnURL string `json:"return_url"`
	CancelURL string `json:"cancel_url"`
}

func (h *Handler) CreatePayment(c *gin.Context) {
	user := middleware.UserFromContext(c)
	if user == nil {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized", nil)
		return
	}
	idempotencyKey := strings.TrimSpace(c.GetHeader("Idempotency-Key"))
	if idempotencyKey == "" {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Idempotency-Key is required", nil)
		return
	}

	var req createPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	output, apiErr := h.payments.Create(paysvc.CreatePaymentInput{
		UserID:         user.ID,
		OrderID:        c.Param("order_id"),
		IdempotencyKey: idempotencyKey,
		Provider:       req.Provider,
		ReturnURL:      req.ReturnURL,
		CancelURL:      req.CancelURL,
	})
	if apiErr != nil {
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message, apiErr.Details)
		return
	}

	response.Created(c, gin.H{
		"payment_id":   output.PaymentID,
		"status":       output.Status,
		"checkout_url": output.CheckoutURL,
	})
}

func (h *Handler) ListPayments(c *gin.Context) {
	user := middleware.UserFromContext(c)
	if user == nil {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized", nil)
		return
	}

	list, apiErr := h.payments.ListByOrderForUser(c.Param("order_id"), user.ID)
	if apiErr != nil {
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message, apiErr.Details)
		return
	}

	items := make([]gin.H, 0, len(list))
	for _, p := range list {
		items = append(items, gin.H{
			"id":                p.ID,
			"txn_type":          p.TxnType,
			"provider":          p.Provider,
			"provider_txn_code": p.ProviderTxnCode,
			"status":            p.Status,
			"amount":            utils.ToAmountString(p.Amount),
			"currency_code":     p.CurrencyCode,
			"created_at":        p.CreatedAt.Format(time.RFC3339Nano),
		})
	}

	response.OK(c, gin.H{"items": items})
}

func (h *Handler) ReceivePaymentWebhook(c *gin.Context) {
	provider := c.Param("provider")
	var payload map[string]any
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	result, apiErr := h.payments.ProcessWebhook(provider, payload)
	if apiErr != nil {
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message, apiErr.Details)
		return
	}

	response.Accepted(c, result)
}
