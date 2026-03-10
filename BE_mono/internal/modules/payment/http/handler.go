package http

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"golf-store/be-mono/internal/modules/payment/dto"
	paysvc "golf-store/be-mono/internal/modules/payment/service"
	entities "golf-store/be-mono/internal/platform/entities"
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

	var req dto.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if !errors.Is(err, io.EOF) {
			response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
			return
		}
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

	response.Created(c, dto.CreatePaymentResponseDTO{
		PaymentID:   output.PaymentID,
		Status:      output.Status,
		CheckoutURL: output.CheckoutURL,
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

	items := make([]dto.PaymentItemDTO, 0, len(list))
	for _, p := range list {
		items = append(items, paymentResponse(p))
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

func paymentResponse(p *entities.PaymentTransaction) dto.PaymentItemDTO {
	if p == nil {
		return dto.PaymentItemDTO{}
	}
	resp := dto.PaymentItemDTO{
		ID:           p.ID,
		TxnType:      p.TxnType,
		Provider:     p.Provider,
		Status:       p.Status,
		Amount:       utils.ToAmountString(p.Amount),
		CurrencyCode: p.CurrencyCode,
		CreatedAt:    p.CreatedAt.Format(time.RFC3339Nano),
	}
	if p.ProviderTxnCode != nil {
		resp.ProviderTxnCode = *p.ProviderTxnCode
	}
	return resp
}
