package repository

import (
	"time"

	"gorm.io/gorm"

	"golf-store/be-mono/internal/platform/db"
)

type ListFilter struct {
	UserID   string
	Status   string
	Page     int
	PageSize int
}

type Repository interface {
	List(filter ListFilter) ([]db.NotificationEntity, int64, error)
	MarkRead(userID string, notificationID string) (int64, error)
	MarkReadAll(userID string) (int64, error)
	FindByID(notificationID string) (*db.NotificationEntity, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGorm(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) List(filter ListFilter) ([]db.NotificationEntity, int64, error) {
	query := r.db.Model(&db.NotificationEntity{}).Where("user_id = ?", filter.UserID)
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PageSize
	entities := make([]db.NotificationEntity, 0)
	if err := query.Order("created_at DESC").Limit(filter.PageSize).Offset(offset).Find(&entities).Error; err != nil {
		return nil, total, err
	}

	return entities, total, nil
}

func (r *GormRepository) MarkRead(userID string, notificationID string) (int64, error) {
	res := r.db.Model(&db.NotificationEntity{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Updates(map[string]any{
			"is_read":    true,
			"updated_at": time.Now().UTC(),
		})
	return res.RowsAffected, res.Error
}

func (r *GormRepository) MarkReadAll(userID string) (int64, error) {
	res := r.db.Model(&db.NotificationEntity{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]any{
			"is_read":    true,
			"updated_at": time.Now().UTC(),
		})
	return res.RowsAffected, res.Error
}

func (r *GormRepository) FindByID(notificationID string) (*db.NotificationEntity, error) {
	var entity db.NotificationEntity
	err := r.db.Where("id = ?", notificationID).Take(&entity).Error
	return &entity, err
}
