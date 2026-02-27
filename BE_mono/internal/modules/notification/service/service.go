package service

import (
	"net/http"
	"strings"
	"time"

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
	UserID   string
	Status   string
	Page     int
	PageSize int
}

type ListOutput struct {
	Items      []*db.NotificationEntity
	Page       int
	PageSize   int
	Total      int
	TotalPages int
}

func (s *Service) List(input ListInput) ListOutput {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 20
	}

	query := s.db.Model(&db.NotificationEntity{}).Where("user_id = ?", input.UserID)
	if strings.TrimSpace(input.Status) != "" {
		query = query.Where("status = ?", strings.ToUpper(strings.TrimSpace(input.Status)))
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return ListOutput{Items: []*db.NotificationEntity{}, Page: input.Page, PageSize: input.PageSize}
	}

	offset := (input.Page - 1) * input.PageSize
	entities := make([]db.NotificationEntity, 0)
	if err := query.Order("created_at DESC").Limit(input.PageSize).Offset(offset).Find(&entities).Error; err != nil {
		return ListOutput{
			Items:      []*db.NotificationEntity{},
			Page:       input.Page,
			PageSize:   input.PageSize,
			Total:      int(total),
			TotalPages: calcTotalPages(int(total), input.PageSize),
		}
	}

	items := make([]*db.NotificationEntity, 0, len(entities))
	for i := range entities {
		items = append(items, cloneNotification(&entities[i]))
	}

	return ListOutput{
		Items:      items,
		Page:       input.Page,
		PageSize:   input.PageSize,
		Total:      int(total),
		TotalPages: calcTotalPages(int(total), input.PageSize),
	}
}

func (s *Service) MarkRead(userID string, notificationID string) (*db.NotificationEntity, *apperrors.APIError) {
	now := time.Now().UTC()
	res := s.db.Model(&db.NotificationEntity{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Updates(map[string]any{
			"is_read":    true,
			"updated_at": now,
		})
	if res.Error != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to update notification"}
	}
	if res.RowsAffected == 0 {
		return nil, apperrors.ErrNotFound
	}

	var entity db.NotificationEntity
	if err := s.db.Where("id = ?", notificationID).Take(&entity).Error; err != nil {
		return nil, apperrors.ErrNotFound
	}
	return cloneNotification(&entity), nil
}

func (s *Service) MarkReadAll(userID string) int {
	res := s.db.Model(&db.NotificationEntity{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]any{
			"is_read":    true,
			"updated_at": time.Now().UTC(),
		})
	if res.Error != nil {
		return 0
	}
	return int(res.RowsAffected)
}

func cloneNotification(entity *db.NotificationEntity) *db.NotificationEntity {
	if entity == nil {
		return nil
	}
	copy := *entity
	copy.CreatedAt = copy.CreatedAt.UTC()
	copy.UpdatedAt = copy.UpdatedAt.UTC()
	if copy.SentAt != nil {
		t := copy.SentAt.UTC()
		copy.SentAt = &t
	}
	return &copy
}

func calcTotalPages(total int, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	return (total + pageSize - 1) / pageSize
}
