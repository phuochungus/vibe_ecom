package service

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"golf-store/be-mono/internal/platform/db"
	apperrors "golf-store/be-mono/internal/shared/errors"
)

type Service struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Service {
	return &Service{db: db}
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

type ListOutput struct {
	Items      []*db.ProductEntity
	Page       int
	PageSize   int
	Total      int
	TotalPages int
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

func (s *Service) List(input ListInput) ListOutput {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 20
	}

	query := s.db.Model(&db.ProductEntity{}).Where("deleted_at IS NULL")

	if !input.AdminView {
		query = query.Where("status = ?", db.ProductStatusActive)
	}
	if strings.TrimSpace(input.Status) != "" {
		query = query.Where("status = ?", strings.ToUpper(strings.TrimSpace(input.Status)))
	}
	if strings.TrimSpace(input.Query) != "" {
		q := "%" + strings.ToLower(strings.TrimSpace(input.Query)) + "%"
		query = query.Where("(LOWER(name) LIKE ? OR LOWER(sku) LIKE ?)", q, q)
	}
	if input.Min != nil {
		query = query.Where("price_cents >= ?", *input.Min)
	}
	if input.Max != nil {
		query = query.Where("price_cents <= ?", *input.Max)
	}

	sortBy := sanitizeSortBy(input.SortBy)
	sortOrder := sanitizeSortOrder(input.SortOrder)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return ListOutput{Items: []*db.ProductEntity{}, Page: input.Page, PageSize: input.PageSize, Total: 0, TotalPages: 0}
	}

	offset := (input.Page - 1) * input.PageSize
	var entities []db.ProductEntity
	if err := query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
		Limit(input.PageSize).
		Offset(offset).
		Find(&entities).Error; err != nil {
		return ListOutput{
			Items:      []*db.ProductEntity{},
			Page:       input.Page,
			PageSize:   input.PageSize,
			Total:      int(total),
			TotalPages: calcTotalPages(int(total), input.PageSize),
		}
	}

	items := make([]*db.ProductEntity, 0, len(entities))
	for i := range entities {
		items = append(items, cloneProduct(&entities[i]))
	}

	return ListOutput{
		Items:      items,
		Page:       input.Page,
		PageSize:   input.PageSize,
		Total:      int(total),
		TotalPages: calcTotalPages(int(total), input.PageSize),
	}
}

func (s *Service) GetByID(productID string, adminView bool) (*db.ProductEntity, *apperrors.APIError) {
	query := s.db.Where("id = ? AND deleted_at IS NULL", productID)
	if !adminView {
		query = query.Where("status = ?", db.ProductStatusActive)
	}

	var entity db.ProductEntity
	err := query.Take(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query product"}
	}

	return cloneProduct(&entity), nil
}

func (s *Service) AdminCreate(input AdminUpsertInput) (*db.ProductEntity, *apperrors.APIError) {
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
		status = db.ProductStatusActive
	}
	entity := &db.ProductEntity{
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
	if err := s.db.Create(entity).Error; err != nil {
		return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: "invalid product payload"}
	}

	return s.GetByID(id, true)
}

func (s *Service) AdminUpdate(productID string, input AdminUpsertInput) (*db.ProductEntity, *apperrors.APIError) {
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

	updateMap["updated_at"] = time.Now().UTC()

	res := s.db.Model(&db.ProductEntity{}).
		Where("id = ? AND deleted_at IS NULL", productID).
		Updates(updateMap)
	if res.Error != nil {
		return nil, &apperrors.APIError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: "invalid product payload"}
	}
	if res.RowsAffected == 0 {
		return nil, apperrors.ErrNotFound
	}

	return s.GetByID(productID, true)
}

func (s *Service) AdminDelete(productID string) *apperrors.APIError {
	now := time.Now().UTC()
	res := s.db.Model(&db.ProductEntity{}).
		Where("id = ? AND deleted_at IS NULL", productID).
		Updates(map[string]any{
			"deleted_at": now,
			"status":     db.ProductStatusDiscontinued,
			"updated_at": now,
		})
	if res.Error != nil {
		return &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to delete product"}
	}
	if res.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func ParseStatus(status string) (string, error) {
	if status == "" {
		return db.ProductStatusActive, nil
	}
	up := strings.ToUpper(strings.TrimSpace(status))
	switch up {
	case db.ProductStatusActive:
		return db.ProductStatusActive, nil
	case db.ProductStatusInactive:
		return db.ProductStatusInactive, nil
	case db.ProductStatusDiscontinued:
		return db.ProductStatusDiscontinued, nil
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

func cloneProduct(entity *db.ProductEntity) *db.ProductEntity {
	if entity == nil {
		return nil
	}
	copy := *entity
	copy.CreatedAt = copy.CreatedAt.UTC()
	copy.UpdatedAt = copy.UpdatedAt.UTC()
	if copy.DeletedAt != nil {
		t := copy.DeletedAt.UTC()
		copy.DeletedAt = &t
	}
	return &copy
}
