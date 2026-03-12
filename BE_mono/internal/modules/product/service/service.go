package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"math"

	"github.com/google/uuid"

	"golf-store/be-mono/internal/modules/product/repository"
	entities "golf-store/be-mono/internal/platform/entities"
	apperrors "golf-store/be-mono/internal/shared/errors"
	"golf-store/be-mono/internal/shared/response"
)

type Service struct {
	repo         repository.Repository
	imageStorage ImageStorage
}

type ImageStorage interface {
	Upload(ctx context.Context, objectKey string, reader io.Reader, size int64, contentType string) (string, error)
}

func New(repo repository.Repository, imageStorage ImageStorage) *Service {
	return &Service{
		repo:         repo,
		imageStorage: imageStorage,
	}
}

type ListInput struct {
	Query     string
	Status    string
	Min       *int64
	Max       *int64
	Page      int
	PageSize  int
	SortBy    string
	SortOrder string
	AdminView bool
}

type AdminUpsertInput struct {
	SKU         string
	Name        string
	Description string
	Price       int64
	Stock       int
	Status      string
	ImageURL    string
}

type AdminUploadImageOutput struct {
	URL         string `json:"url"`
	ObjectKey   string `json:"object_key"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}

const MaxProductImageUploadBytes = 10 << 20

var productImageNameSanitizer = regexp.MustCompile(`[^a-z0-9]+`)

func (s *Service) List(input ListInput) response.PageDto[*entities.Product] {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 20
	}

	filter := repository.ListFilter{
		Query:     input.Query,
		Status:    input.Status,
		Min:       input.Min,
		Max:       input.Max,
		Page:      input.Page,
		PageSize:  input.PageSize,
		SortBy:    input.SortBy,
		SortOrder: input.SortOrder,
		AdminView: input.AdminView,
	}

	items, total, err := s.repo.List(filter)

	if err != nil {
		return response.PageDto[*entities.Product]{
			Items: []*entities.Product{},
			Pagination: response.PageMeta{
				Page:       input.Page,
				PageSize:   input.PageSize,
				Total:      0,
				TotalPages: 0,
			},
		}
	}
	if input.AdminView == false {
		for i := range items {
			items[i].Stock = math.MaxInt
		}
	}

	return response.PageDto[*entities.Product]{
		Items: items,
		Pagination: response.PageMeta{
			Page:       input.Page,
			PageSize:   input.PageSize,
			Total:      int(total),
			TotalPages: calcTotalPages(int(total), input.PageSize),
		},
	}
}

func (s *Service) GetByID(productID string, adminView bool) (*entities.Product, *apperrors.APIError) {
	entity, err := s.repo.FindByID(productID, adminView)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, apperrors.ErrNotFound
		}
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query product"}
	}

	return entity, nil
}

func (s *Service) AdminCreate(input AdminUpsertInput) (*entities.Product, *apperrors.APIError) {
	now := time.Now().UTC()
	if input.Price <= 0 {
		return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: "price must be greater than 0"}
	}
	if input.Stock < 0 {
		return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: "stock cannot be negative"}
	}

	id := uuid.NewString()
	status := input.Status
	if status == "" {
		status = entities.ProductStatusActive
	}
	entity := &entities.Product{
		ID:          id,
		SKU:         strings.TrimSpace(input.SKU),
		Name:        strings.TrimSpace(input.Name),
		Description: strings.TrimSpace(input.Description),
		Price:       input.Price,
		Stock:       input.Stock,
		Status:      status,
		ImageURL:    strings.TrimSpace(input.ImageURL),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.repo.Create(entity); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: "invalid product payload"}
	}

	return s.GetByID(id, true)
}

func (s *Service) AdminUpdate(productID string, input AdminUpsertInput) (*entities.Product, *apperrors.APIError) {
	updateMap := map[string]any{}

	if strings.TrimSpace(input.SKU) != "" {
		updateMap["sku"] = strings.TrimSpace(input.SKU)
	}
	if strings.TrimSpace(input.Name) != "" {
		updateMap["name"] = strings.TrimSpace(input.Name)
	}
	if strings.TrimSpace(input.Description) != "" {
		updateMap["description"] = strings.TrimSpace(input.Description)
	}
	if input.Price > 0 {
		updateMap["price"] = input.Price
	}
	if input.Stock >= 0 {
		updateMap["stock"] = input.Stock
	}
	if input.Status != "" {
		updateMap["status"] = input.Status
	}
	if strings.TrimSpace(input.ImageURL) != "" {
		updateMap["image_url"] = strings.TrimSpace(input.ImageURL)
	}
	if len(updateMap) == 0 {
		return s.GetByID(productID, true)
	}

	rowsAffected, err := s.repo.Update(productID, updateMap)
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: "invalid product payload"}
	}
	if rowsAffected == 0 {
		return nil, apperrors.ErrNotFound
	}

	return s.GetByID(productID, true)
}

func (s *Service) AdminDelete(productID string) *apperrors.APIError {
	now := time.Now().UTC()
	updates := map[string]any{
		"deleted_at": now,
		"status":     entities.ProductStatusDiscontinued,
		"updated_at": now,
	}

	rowsAffected, err := s.repo.SoftDelete(productID, updates)
	if err != nil {
		return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to delete product"}
	}
	if rowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (s *Service) AdminUploadImage(ctx context.Context, fileName string, content []byte) (*AdminUploadImageOutput, *apperrors.APIError) {
	if s.imageStorage == nil {
		return nil, &apperrors.APIError{
			Status:  http.StatusServiceUnavailable,
			Code:    "IMAGE_STORAGE_NOT_CONFIGURED",
			Message: "Image storage is not configured",
		}
	}

	if len(content) == 0 {
		return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: "image file is required"}
	}
	if len(content) > MaxProductImageUploadBytes {
		return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: "image must be 10MB or smaller"}
	}

	contentType := http.DetectContentType(content)
	if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(contentType)), "image/") {
		return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: "uploaded file must be an image"}
	}

	objectKey := buildProductImageObjectKey(fileName, contentType)
	url, err := s.imageStorage.Upload(ctx, objectKey, bytes.NewReader(content), int64(len(content)), contentType)
	if err != nil {
		return nil, &apperrors.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "IMAGE_UPLOAD_FAILED",
			Message: "Failed to upload image",
			Details: err.Error(),
		}
	}

	return &AdminUploadImageOutput{
		URL:         url,
		ObjectKey:   objectKey,
		ContentType: contentType,
		Size:        int64(len(content)),
	}, nil
}

func ParseStatus(status string) (string, error) {
	if status == "" {
		return entities.ProductStatusActive, nil
	}
	up := strings.ToUpper(strings.TrimSpace(status))
	switch up {
	case entities.ProductStatusActive:
		return entities.ProductStatusActive, nil
	case entities.ProductStatusInactive:
		return entities.ProductStatusInactive, nil
	case entities.ProductStatusDiscontinued:
		return entities.ProductStatusDiscontinued, nil
	default:
		return "", fmt.Errorf("invalid product status")
	}
}

func sanitizeSortBy(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "price":
		return "price"
	case "created_at":
		return "created_at"
	default:
		return "name"
	}
}

func sanitizeSortOrder(value string) string {
	if strings.EqualFold(strings.TrimSpace(value), "desc") {
		return "DESC"
	}
	return "ASC"
}

func calcTotalPages(total int, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	return (total + pageSize - 1) / pageSize
}

func buildProductImageObjectKey(fileName string, contentType string) string {
	ext := imageExtension(fileName, contentType)
	baseName := sanitizeImageBaseName(strings.TrimSuffix(filepath.Base(strings.TrimSpace(fileName)), filepath.Ext(strings.TrimSpace(fileName))))
	if baseName == "" {
		baseName = "product-image"
	}

	return fmt.Sprintf(
		"products/%s/%s-%s%s",
		time.Now().UTC().Format("2006/01"),
		uuid.NewString(),
		baseName,
		ext,
	)
}

func imageExtension(fileName string, contentType string) string {
	switch strings.ToLower(strings.TrimSpace(contentType)) {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	case "image/avif":
		return ".avif"
	}

	ext := strings.ToLower(strings.TrimSpace(filepath.Ext(fileName)))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif", ".avif":
		if ext == ".jpeg" {
			return ".jpg"
		}
		return ext
	default:
		return ".bin"
	}
}

func sanitizeImageBaseName(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	normalized = productImageNameSanitizer.ReplaceAllString(normalized, "-")
	return strings.Trim(normalized, "-")
}
