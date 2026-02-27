package repository

import (
	"fmt"
	"strings"

	"gorm.io/gorm"

	"golf-store/be-mono/internal/platform/db"
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
	List(filter ListFilter) ([]*db.ProductEntity, int64, error)
	FindByID(id string, adminView bool) (*db.ProductEntity, error)
	Create(entity *db.ProductEntity) error
	Update(id string, updates map[string]any) (int64, error)
	SoftDelete(id string, updates map[string]any) (int64, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGorm(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) List(filter ListFilter) ([]*db.ProductEntity, int64, error) {
	query := r.db.Model(&db.ProductEntity{}).Where("deleted_at IS NULL")

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
	var entities []db.ProductEntity
	if err := query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
		Limit(filter.PageSize).
		Offset(offset).
		Find(&entities).Error; err != nil {
		return nil, total, err
	}

	items := make([]*db.ProductEntity, 0, len(entities))
	for i := range entities {
		items = append(items, &entities[i])
	}

	return items, total, nil
}

func (r *GormRepository) FindByID(id string, adminView bool) (*db.ProductEntity, error) {
	query := r.db.Where("id = ? AND deleted_at IS NULL", id)
	if !adminView {
		query = query.Where("status = ?", db.ProductStatusActive)
	}

	var entity db.ProductEntity
	err := query.Take(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *GormRepository) Create(entity *db.ProductEntity) error {
	return r.db.Create(entity).Error
}

func (r *GormRepository) Update(id string, updates map[string]any) (int64, error) {
	res := r.db.Model(&db.ProductEntity{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates)
	return res.RowsAffected, res.Error
}

func (r *GormRepository) SoftDelete(id string, updates map[string]any) (int64, error) {
	res := r.db.Model(&db.ProductEntity{}).
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
