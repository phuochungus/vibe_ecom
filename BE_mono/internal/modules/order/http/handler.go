package http

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"golf-store/be-mono/internal/modules/order/dto"
	ordersvc "golf-store/be-mono/internal/modules/order/service"
	entities "golf-store/be-mono/internal/platform/entities"
	"golf-store/be-mono/internal/shared/middleware"
	"golf-store/be-mono/internal/shared/response"
	"golf-store/be-mono/internal/shared/utils"
)

type Handler struct {
	orders *ordersvc.Service
}

func New(orders *ordersvc.Service) *Handler {
	return &Handler{orders: orders}
}

func (h *Handler) RegisterUser(rg *gin.RouterGroup) {
	rg.POST("/orders", h.CreateOrder)
	rg.GET("/orders", h.ListOrders)
	rg.GET("/orders/:order_id", h.GetOrder)
	rg.POST("/orders/:order_id/cancel", h.CancelOrder)
	rg.GET("/orders/:order_id/tracking", h.GetTracking)
}

func (h *Handler) RegisterAdmin(rg *gin.RouterGroup) {
	rg.GET("/orders", h.AdminListOrders)
	rg.GET("/orders/:order_id", h.AdminGetOrder)
	rg.PATCH("/orders/:order_id/status", h.AdminUpdateOrderStatus)
}

func (h *Handler) CreateOrder(c *gin.Context) {
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

	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	items := make([]dto.CreateOrderItemInput, 0, len(req.Items))
	for _, i := range req.Items {
		items = append(items, dto.CreateOrderItemInput{ProductID: i.ProductID, Quantity: i.Quantity})
	}

	order, apiErr := h.orders.Create(dto.CreateOrderInput{
		UserID:         user.ID,
		IdempotencyKey: idempotencyKey,
		Items:          items,
		Shipping: dto.ShippingAddress{
			RecipientName:  req.ShippingAddress.RecipientName,
			RecipientPhone: req.ShippingAddress.RecipientPhone,
			Line1:          req.ShippingAddress.Line1,
			Line2:          req.ShippingAddress.Line2,
			Ward:           req.ShippingAddress.Ward,
			District:       req.ShippingAddress.District,
			City:           req.ShippingAddress.City,
			Province:       req.ShippingAddress.Province,
			PostalCode:     req.ShippingAddress.PostalCode,
			CountryCode:    req.ShippingAddress.CountryCode,
		},
		CustomerNote: req.CustomerNote,
	})
	if apiErr != nil {
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message, apiErr.Details)
		return
	}

	response.Created(c, orderSummaryResponse(order))
}

func (h *Handler) ListOrders(c *gin.Context) {
	user := middleware.UserFromContext(c)
	if user == nil {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized", nil)
		return
	}

	from, _ := parseOptionalTime(c.Query("from"))
	to, _ := parseOptionalTime(c.Query("to"))
	page := parseIntDefault(c.Query("page"), 1)
	pageSize := parseIntDefault(c.Query("page_size"), 20)

	out := h.orders.List(dto.ListInput{
		UserID:   user.ID,
		Status:   c.Query("status"),
		From:     from,
		To:       to,
		Page:     page,
		PageSize: pageSize,
		Admin:    false,
	})

	items := make([]dto.OrderSummaryDTO, 0, len(out.Items))
	for _, o := range out.Items {
		items = append(items, orderSummaryResponse(o))
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

func (h *Handler) GetOrder(c *gin.Context) {
	user := middleware.UserFromContext(c)
	order, apiErr := h.orders.GetByIDForUser(c.Param("order_id"), user.ID)
	if apiErr != nil {
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message, apiErr.Details)
		return
	}
	response.OK(c, orderDetailResponse(order))
}

func (h *Handler) CancelOrder(c *gin.Context) {
	user := middleware.UserFromContext(c)
	var req dto.CancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	order, apiErr := h.orders.CancelByUser(c.Param("order_id"), user.ID, req.Reason)
	if apiErr != nil {
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message, apiErr.Details)
		return
	}
	response.OK(c, orderSummaryResponse(order))
}

func (h *Handler) GetTracking(c *gin.Context) {
	user := middleware.UserFromContext(c)
	timeline, current, apiErr := h.orders.TrackingForUser(c.Param("order_id"), user.ID)
	if apiErr != nil {
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message, apiErr.Details)
		return
	}
	items := make([]gin.H, 0, len(timeline))
	for _, e := range timeline {
		description := ""
		if e.Description != nil {
			description = *e.Description
		}
		items = append(items, gin.H{
			"status":      e.ToStatus,
			"source_type": e.SourceType,
			"occurred_at": e.OccurredAt.Format(time.RFC3339Nano),
			"description": description,
		})
	}
	response.OK(c, gin.H{
		"order_id":       c.Param("order_id"),
		"current_status": current,
		"timeline":       items,
	})
}

func (h *Handler) AdminListOrders(c *gin.Context) {
	from, _ := parseOptionalTime(c.Query("from"))
	to, _ := parseOptionalTime(c.Query("to"))
	page := parseIntDefault(c.Query("page"), 1)
	pageSize := parseIntDefault(c.Query("page_size"), 20)

	out := h.orders.List(dto.ListInput{
		Status:   c.Query("status"),
		From:     from,
		To:       to,
		Page:     page,
		PageSize: pageSize,
		Admin:    true,
	})

	items := make([]dto.OrderSummaryDTO, 0, len(out.Items))
	for _, o := range out.Items {
		items = append(items, orderSummaryResponse(o))
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

func (h *Handler) AdminGetOrder(c *gin.Context) {
	order, apiErr := h.orders.GetByIDAdmin(c.Param("order_id"))
	if apiErr != nil {
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message, apiErr.Details)
		return
	}
	response.OK(c, orderDetailResponse(order))
}

func (h *Handler) AdminUpdateOrderStatus(c *gin.Context) {
	var req dto.AdminUpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}
	toStatus, err := ordersvc.ParseOrderStatus(req.ToStatus)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "invalid to_status", nil)
		return
	}

	order, apiErr := h.orders.AdminUpdateStatus(c.Param("order_id"), toStatus, req.Reason)
	if apiErr != nil {
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message, apiErr.Details)
		return
	}
	response.OK(c, orderSummaryResponse(order))
}

func orderSummaryResponse(o *entities.Order) dto.OrderSummaryDTO {
	if o == nil {
		return dto.OrderSummaryDTO{}
	}
	result := dto.OrderSummaryDTO{
		ID:             o.ID,
		OrderCode:      o.OrderCode,
		OrderStatus:    o.OrderStatus,
		PaymentStatus:  o.PaymentStatus,
		SubtotalAmount: utils.ToAmountString(o.SubtotalAmount),
		DiscountAmount: utils.ToAmountString(o.DiscountAmount),
		ShippingFee:    utils.ToAmountString(o.ShippingFee),
		TotalAmount:    utils.ToAmountString(o.TotalAmount),
		PlacedAt:       o.PlacedAt.Format(time.RFC3339Nano),
	}
	if o.PaymentDueAt != nil {
		result.PaymentDueAt = o.PaymentDueAt.Format(time.RFC3339Nano)
	}
	return result
}

func orderDetailResponse(o *entities.Order) dto.OrderDetailDTO {
	if o == nil || o.Items == nil {
		return dto.OrderDetailDTO{}
	}
	items := make([]dto.OrderItemDTO, 0, len(o.Items))
	for _, item := range o.Items {
		items = append(items, dto.OrderItemDTO{
			ProductID: item.ProductID,
			SKU:       item.SKU,
			Name:      item.Name,
			UnitPrice: utils.ToAmountString(item.UnitPrice),
			Quantity:  item.Quantity,
			LineTotal: utils.ToAmountString(item.LineTotal),
		})
	}
	return dto.OrderDetailDTO{
		OrderSummaryDTO: orderSummaryResponse(o),
		Items:           items,
		Payments:        []any{},
	}
}

func parseOptionalTime(value string) (*time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, err
	}
	utc := t.UTC()
	return &utc, nil
}

func parseIntDefault(value string, fallback int) int {
	v, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || v <= 0 {
		return fallback
	}
	return v
}
