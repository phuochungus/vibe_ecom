package repository

import (
	"fmt"
	"strings"

	"gorm.io/gorm"

	"golf-store/be-mono/internal/platform/db"
	entities "golf-store/be-mono/internal/platform/entities"
)

type ListFilter struct {
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

type Repository interface {
	List(filter ListFilter) ([]*entities.Product, int64, error)
	FindByID(id string, adminView bool) (*entities.Product, error)
	Create(entity *entities.Product) error
	Update(id string, updates map[string]any) (int64, error)
	SoftDelete(id string, updates map[string]any) (int64, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGorm(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) List(filter ListFilter) ([]*entities.Product, int64, error) {
	query := r.db.Model(&entities.Product{}).Where("deleted_at IS NULL")

	if !filter.AdminView {
		query = query.Where("status = ?", db.ProductStatusActive)
	}
	if strings.TrimSpace(filter.Status) != "" {
		query = query.Where("status = ?", strings.ToUpper(strings.TrimSpace(filter.Status)))
	}
	if strings.TrimSpace(filter.Query) != "" {
		q := "%" + strings.ToLower(strings.TrimSpace(filter.Query)) + "%"
		query = query.Where("(LOWER(name) LIKE ? OR LOWER(sku) LIKE ?)", q, q)
	}
	if filter.Min != nil {
		query = query.Where("price_cents >= ?", *filter.Min)
	}
	if filter.Max != nil {
		query = query.Where("price_cents <= ?", *filter.Max)
	}

	sortBy := sanitizeSortBy(filter.SortBy)
	sortOrder := sanitizeSortOrder(filter.SortOrder)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PageSize
	var products []entities.Product
	if err := query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
		Limit(filter.PageSize).
		Offset(offset).
		Find(&products).Error; err != nil {
		return nil, total, err
	}

	items := make([]*entities.Product, 0, len(products))
	for i := range products {
		items = append(items, &products[i])
	}

	return items, total, nil
}

func (r *GormRepository) FindByID(id string, adminView bool) (*entities.Product, error) {
	query := r.db.Where("id = ? AND deleted_at IS NULL", id)
	if !adminView {
		query = query.Where("status = ?", db.ProductStatusActive)
	}

	var entity entities.Product
	err := query.Take(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *GormRepository) Create(entity *entities.Product) error {
	return r.db.Create(entity).Error
}

func (r *GormRepository) Update(id string, updates map[string]any) (int64, error) {
	res := r.db.Model(&entities.Product{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates)
	return res.RowsAffected, res.Error
}

func (r *GormRepository) SoftDelete(id string, updates map[string]any) (int64, error) {
	res := r.db.Model(&entities.Product{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates)
	return res.RowsAffected, res.Error
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
