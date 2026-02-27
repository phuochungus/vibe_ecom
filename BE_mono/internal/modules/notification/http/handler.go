package http

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	notisvc "golf-store/be-mono/internal/modules/notification/service"
	"golf-store/be-mono/internal/platform/db"
	"golf-store/be-mono/internal/shared/middleware"
	"golf-store/be-mono/internal/shared/response"
)

type Handler struct {
	notifications *notisvc.Service
}

func New(notifications *notisvc.Service) *Handler {
	return &Handler{notifications: notifications}
}

func (h *Handler) RegisterUser(rg *gin.RouterGroup) {
	rg.GET("/notifications", h.ListNotifications)
	rg.PATCH("/notifications/:notification_id/read", h.MarkRead)
	rg.PATCH("/notifications/read-all", h.MarkReadAll)
}

type notificationResponseDTO struct {
	ID        string `json:"id"`
	Channel   string `json:"channel"`
	EventType string `json:"event_type"`
	EventKey  string `json:"event_key"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Status    string `json:"status"`
	Read      bool   `json:"read"`
	CreatedAt string `json:"created_at"`
	SentAt    string `json:"sent_at,omitempty"`
}

func (h *Handler) ListNotifications(c *gin.Context) {
	user := middleware.UserFromContext(c)
	if user == nil {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized", nil)
		return
	}

	out := h.notifications.List(notisvc.ListInput{
		UserID:   user.ID,
		Status:   c.Query("status"),
		Page:     parseIntDefault(c.Query("page"), 1),
		PageSize: parseIntDefault(c.Query("page_size"), 20),
	})

	items := make([]notificationResponseDTO, 0, len(out.Items))
	for _, n := range out.Items {
		items = append(items, notificationResponse(n))
	}

	response.OK(c, gin.H{
		"items": items,
		"pagination": gin.H{
			"page":        out.Page,
			"page_size":   out.PageSize,
			"total":       out.Total,
			"total_pages": out.TotalPages,
		},
	})
}

func (h *Handler) MarkRead(c *gin.Context) {
	user := middleware.UserFromContext(c)
	n, apiErr := h.notifications.MarkRead(user.ID, c.Param("notification_id"))
	if apiErr != nil {
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message, apiErr.Details)
		return
	}
	response.OK(c, notificationResponse(n))
}

func (h *Handler) MarkReadAll(c *gin.Context) {
	user := middleware.UserFromContext(c)
	count := h.notifications.MarkReadAll(user.ID)
	response.OK(c, gin.H{"updated_count": count})
}

func notificationResponse(n *db.NotificationEntity) notificationResponseDTO {
	if n == nil {
		return notificationResponseDTO{}
	}
	resp := notificationResponseDTO{
		ID:        n.ID,
		Channel:   n.Channel,
		EventType: n.EventType,
		EventKey:  n.EventKey,
		Title:     n.Title,
		Content:   n.Content,
		Status:    n.Status,
		Read:      n.IsRead,
		CreatedAt: n.CreatedAt.Format(time.RFC3339Nano),
	}
	if n.SentAt != nil {
		resp.SentAt = n.SentAt.Format(time.RFC3339Nano)
	}
	return resp
}

func parseIntDefault(value string, fallback int) int {
	v, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || v <= 0 {
		return fallback
	}
	return v
}
