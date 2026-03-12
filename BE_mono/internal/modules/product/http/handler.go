package http

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"golf-store/be-mono/internal/modules/product/dto"
	prodsvc "golf-store/be-mono/internal/modules/product/service"
	"golf-store/be-mono/internal/shared/response"
)

type Handler struct {
	products *prodsvc.Service
}

func New(products *prodsvc.Service) *Handler {
	return &Handler{products: products}
}

func (h *Handler) RegisterPublic(rg *gin.RouterGroup) {
	rg.GET("/products", h.ListProducts)
	rg.GET("/products/:product_id", h.GetProduct)
}

func (h *Handler) RegisterAdmin(rg *gin.RouterGroup) {
	rg.GET("/products", h.AdminListProducts)
	rg.GET("/products/:product_id", h.AdminGetProduct)
	rg.POST("/products/upload-image", h.AdminUploadProductImage)
	rg.POST("/products", h.AdminCreateProduct)
	rg.PATCH("/products/:product_id", h.AdminUpdateProduct)
	rg.DELETE("/products/:product_id", h.AdminDeleteProduct)
}

func (h *Handler) ListProducts(c *gin.Context) {
	//nullable int64
	min := parseIntPtr(c.Query("min_price"))
	max := parseIntPtr(c.Query("max_price"))

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

	response.OK(c, out)
}

func (h *Handler) GetProduct(c *gin.Context) {
	productID := c.Param("product_id")
	p, apiErr := h.products.GetByID(productID, false)
	if apiErr != nil {
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message, apiErr.Details)
		return
	}
	response.OK(c, p)
}

func (h *Handler) AdminListProducts(c *gin.Context) {
	min := parseIntPtr(c.Query("min_price"))
	max := parseIntPtr(c.Query("max_price"))

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
		AdminView: true,
	})

	response.OK(c, out)
}

func (h *Handler) AdminGetProduct(c *gin.Context) {
	productID := c.Param("product_id")
	p, apiErr := h.products.GetByID(productID, true)
	if apiErr != nil {
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message, apiErr.Details)
		return
	}
	response.OK(c, p)
}

func (h *Handler) AdminCreateProduct(c *gin.Context) {
	var req dto.AdminUpsertRequest
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
	response.Created(c, product)
}

func (h *Handler) AdminUpdateProduct(c *gin.Context) {
	productID := c.Param("product_id")
	var req dto.AdminUpsertRequest
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

	response.OK(c, product)
}

func (h *Handler) AdminDeleteProduct(c *gin.Context) {
	apiErr := h.products.AdminDelete(c.Param("product_id"))
	if apiErr != nil {
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message, apiErr.Details)
		return
	}
	response.NoContent(c)
}

func (h *Handler) AdminUploadProductImage(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "image file is required", nil)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "failed to read uploaded file", nil)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	content, err := io.ReadAll(io.LimitReader(file, prodsvc.MaxProductImageUploadBytes+1))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "failed to read uploaded file", nil)
		return
	}
	if len(content) > prodsvc.MaxProductImageUploadBytes {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "image must be 10MB or smaller", nil)
		return
	}

	result, apiErr := h.products.AdminUploadImage(c.Request.Context(), fileHeader.Filename, content)
	if apiErr != nil {
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message, apiErr.Details)
		return
	}

	response.OK(c, result)
}

func parseIntDefault(value string, fallback int) int {
	v, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || v <= 0 {
		return fallback
	}
	return v
}

func parseIntPtr(value string) *int64 {
	if value == "" {
		return nil
	}
	v, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err != nil || v <= 0 {
		return nil
	}
	return &v
}
