package service

import (
	"net/http"
	"strings"

	"golf-store/be-mono/internal/modules/notification/repository"
	entities "golf-store/be-mono/internal/platform/entities"
	apperrors "golf-store/be-mono/internal/shared/errors"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

type ListInput struct {
	UserID   string
	Status   string
	Page     int
	PageSize int
}

type ListOutput struct {
	Items      []*entities.Notification
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

	filter := repository.ListFilter{
		UserID:   input.UserID,
		Status:   strings.ToUpper(strings.TrimSpace(input.Status)),
		Page:     input.Page,
		PageSize: input.PageSize,
	}

	results, total, err := s.repo.List(filter)
	if err != nil {
		return ListOutput{
			Items:      []*entities.Notification{},
			Page:       input.Page,
			PageSize:   input.PageSize,
			Total:      0,
			TotalPages: 0,
		}
	}

	items := make([]*entities.Notification, 0, len(results))
	for i := range results {
		items = append(items, results[i])
	}

	return ListOutput{
		Items:      items,
		Page:       input.Page,
		PageSize:   input.PageSize,
		Total:      int(total),
		TotalPages: calcTotalPages(int(total), input.PageSize),
	}
}

func (s *Service) MarkRead(userID string, notificationID string) (*entities.Notification, *apperrors.APIError) {
	rows, err := s.repo.MarkRead(userID, notificationID)
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to update notification"}
	}
	if rows == 0 {
		return nil, apperrors.ErrNotFound
	}

	entity, err := s.repo.FindByID(notificationID)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.ErrNotFound // For any other fetch error post-update, returning not found is safer
	}
	return entity, nil
}

func (s *Service) MarkReadAll(userID string) int {
	rows, _ := s.repo.MarkReadAll(userID)
	return int(rows)
}

func calcTotalPages(total int, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	return (total + pageSize - 1) / pageSize
}
