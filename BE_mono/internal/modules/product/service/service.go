package service

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"golf-store/be-mono/internal/modules/product/repository"
	entities "golf-store/be-mono/internal/platform/entities"
	apperrors "golf-store/be-mono/internal/shared/errors"
	"golf-store/be-mono/internal/shared/response"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
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
			items[i].Stock = 0
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
		updateMap["price_cents"] = input.Price
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
		return "price_cents"
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
