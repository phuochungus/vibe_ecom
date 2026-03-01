package dto

type NotificationResponseDTO struct {
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
