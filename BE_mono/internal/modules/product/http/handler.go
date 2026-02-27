package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	prodsvc "golf-store/be-mono/internal/modules/product/service"
	"golf-store/be-mono/internal/platform/db"
	"golf-store/be-mono/internal/shared/response"
	"golf-store/be-mono/internal/shared/utils"
)

type Handler struct {
	products *prodsvc.Service
}

func New(products *prodsvc.Service) *Handler {
	return &Handler{products: products}
}

type productResponseDTO struct {
	ID          string `json:"id"`
	SKU         string `json:"sku"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Stock       int    `json:"stock"`
	Status      string `json:"status"`
	ImageURL    string `json:"image_url"`
}

func (h *Handler) RegisterPublic(rg *gin.RouterGroup) {
	rg.GET("/products", h.ListProducts)
	rg.GET("/products/:product_id", h.GetProduct)
}

func (h *Handler) RegisterAdmin(rg *gin.RouterGroup) {
	rg.POST("/products", h.AdminCreateProduct)
	rg.PATCH("/products/:product_id", h.AdminUpdateProduct)
	rg.DELETE("/products/:product_id", h.AdminDeleteProduct)
}

func (h *Handler) ListProducts(c *gin.Context) {
	//nullable int64
	var min, max *int64
	if c.Query("min_price") != "" {
		val, err := strconv.ParseInt(c.Query("min_price"), 10, 64)
		if err == nil {
			min = &val
		}
	}
	if c.Query("max_price") != "" {
		val, err := strconv.ParseInt(c.Query("max_price"), 10, 64)
		if err == nil {
			max = &val
		}
	}

	page := parseIntDefault(c.Query("page"), 1)
	pageSize := parseIntDefault(c.Query("page_size"), 20)

	out := h.products.List(prodsvc.ListInput{
		Query:     c.Query("q"),
		Status:    c.Query("status"),
		Min:       min,
		Max:       max,
		Page:      page,
		PageSize:  pageSize,
		SortBy:    c.Query("sort"),
		SortOrder: c.Query("order"),
		AdminView: false,
	})

	items := make([]productResponseDTO, 0, len(out.Items))
	for _, p := range out.Items {
		items = append(items, productResponse(p))
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

func (h *Handler) GetProduct(c *gin.Context) {
	productID := c.Param("product_id")
	p, apiErr := h.products.GetByID(productID, false)
	if apiErr != nil {
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message, apiErr.Details)
		return
	}
	response.OK(c, productResponse(p))
}

type adminUpsertRequest struct {
	SKU         string `json:"sku"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Stock       *int   `json:"stock"`
	Status      string `json:"status"`
	ImageURL    string `json:"image_url"`
}

func (h *Handler) AdminCreateProduct(c *gin.Context) {
	var req adminUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	priceCents, err := strconv.ParseInt(req.Price, 10, 64)
	if err != nil || priceCents <= 0 {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "price must be greater than 0", nil)
		return
	}
	status, err := prodsvc.ParseStatus(req.Status)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "invalid product status", nil)
		return
	}
	if req.Stock == nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "stock is required", nil)
		return
	}

	product, apiErr := h.products.AdminCreate(prodsvc.AdminUpsertInput{
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		Price:       priceCents,
		Stock:       *req.Stock,
		Status:      status,
		ImageURL:    req.ImageURL,
	})
	if apiErr != nil {
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message, apiErr.Details)
		return
	}
	response.Created(c, productResponse(product))
}

func (h *Handler) AdminUpdateProduct(c *gin.Context) {
	productID := c.Param("product_id")
	var req adminUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	var priceCents int64
	if strings.TrimSpace(req.Price) != "" {
		parsed, err := strconv.ParseInt(req.Price, 10, 64)
		if err != nil || parsed <= 0 {
			response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "price must be greater than 0", nil)
			return
		}
		priceCents = parsed
	}

	status := ""
	if strings.TrimSpace(req.Status) != "" {
		parsedStatus, err := prodsvc.ParseStatus(req.Status)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "invalid product status", nil)
			return
		}
		status = parsedStatus
	}
	stock := -1
	if req.Stock != nil {
		stock = *req.Stock
	}

	product, apiErr := h.products.AdminUpdate(productID, prodsvc.AdminUpsertInput{
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		Price:       priceCents,
		Stock:       stock,
		Status:      status,
		ImageURL:    req.ImageURL,
	})
	if apiErr != nil {
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message, apiErr.Details)
		return
	}

	response.OK(c, productResponse(product))
}

func (h *Handler) AdminDeleteProduct(c *gin.Context) {
	apiErr := h.products.AdminDelete(c.Param("product_id"))
	if apiErr != nil {
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message, apiErr.Details)
		return
	}
	response.NoContent(c)
}

func productResponse(p *db.ProductEntity) productResponseDTO {
	if p == nil {
		return productResponseDTO{}
	}
	return productResponseDTO{
		ID:          p.ID,
		SKU:         p.SKU,
		Name:        p.Name,
		Description: p.Description,
		Price:       utils.ToAmountString(p.Price),
		Stock:       p.Stock,
		Status:      p.Status,
		ImageURL:    p.ImageURL,
	}
}

func parseIntDefault(value string, fallback int) int {
	v, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || v <= 0 {
		return fallback
	}
	return v
}
