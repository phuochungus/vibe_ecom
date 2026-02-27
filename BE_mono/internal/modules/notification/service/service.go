package service

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	apperrors "golf-store/be-mono/internal/shared/errors"
	"golf-store/be-mono/internal/shared/model"
)

type Service struct {
	db *sql.DB
}

func New(db *sql.DB) *Service {
	return &Service{db: db}
}

type ListInput struct {
	UserID   string
	Status   string
	Page     int
	PageSize int
}

type ListOutput struct {
	Items      []*model.Notification
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

	where := "user_id = ?"
	args := []any{input.UserID}
	if strings.TrimSpace(input.Status) != "" {
		where += " AND status = ?"
		args = append(args, strings.ToUpper(strings.TrimSpace(input.Status)))
	}

	var total int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM notifications WHERE `+where, args...).Scan(&total); err != nil {
		return ListOutput{Items: []*model.Notification{}, Page: input.Page, PageSize: input.PageSize}
	}

	offset := (input.Page - 1) * input.PageSize
	listArgs := append(args, input.PageSize, offset)
	rows, err := s.db.Query(
		`SELECT id, user_id, channel, event_type, event_key, title, content, status, is_read, sent_at, created_at, updated_at
		   FROM notifications
		  WHERE `+where+`
		  ORDER BY created_at DESC
		  LIMIT ? OFFSET ?`,
		listArgs...,
	)
	if err != nil {
		return ListOutput{Items: []*model.Notification{}, Page: input.Page, PageSize: input.PageSize, Total: total, TotalPages: calcTotalPages(total, input.PageSize)}
	}
	defer rows.Close()

	items := make([]*model.Notification, 0)
	for rows.Next() {
		n, err := scanNotification(rows)
		if err != nil {
			continue
		}
		items = append(items, n)
	}

	return ListOutput{
		Items:      items,
		Page:       input.Page,
		PageSize:   input.PageSize,
		Total:      total,
		TotalPages: calcTotalPages(total, input.PageSize),
	}
}

func (s *Service) MarkRead(userID string, notificationID string) (*model.Notification, *apperrors.APIError) {
	now := time.Now().UTC()
	res, err := s.db.Exec(
		`UPDATE notifications
		    SET is_read = 1, updated_at = ?
		  WHERE id = ? AND user_id = ?`,
		now, notificationID, userID,
	)
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to update notification"}
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return nil, apperrors.ErrNotFound
	}

	row := s.db.QueryRow(
		`SELECT id, user_id, channel, event_type, event_key, title, content, status, is_read, sent_at, created_at, updated_at
		   FROM notifications
		  WHERE id = ?`,
		notificationID,
	)
	n, err := scanNotification(row)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	return n, nil
}

func (s *Service) MarkReadAll(userID string) int {
	res, err := s.db.Exec(
		`UPDATE notifications
		    SET is_read = 1, updated_at = UTC_TIMESTAMP(3)
		  WHERE user_id = ? AND is_read = 0`,
		userID,
	)
	if err != nil {
		return 0
	}
	affected, _ := res.RowsAffected()
	return int(affected)
}

func scanNotification(scanner interface {
	Scan(dest ...any) error
}) (*model.Notification, error) {
	notification := &model.Notification{}
	var status string
	var isRead bool
	var sentAt sql.NullTime
	if err := scanner.Scan(
		&notification.ID,
		&notification.UserID,
		&notification.Channel,
		&notification.EventType,
		&notification.EventKey,
		&notification.Title,
		&notification.Content,
		&status,
		&isRead,
		&sentAt,
		&notification.CreatedAt,
		&notification.UpdatedAt,
	); err != nil {
		return nil, err
	}
	notification.Status = model.NotificationStatus(status)
	notification.Read = isRead
	if sentAt.Valid {
		t := sentAt.Time.UTC()
		notification.SentAt = &t
	}
	return notification, nil
}

func calcTotalPages(total int, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	return (total + pageSize - 1) / pageSize
}
