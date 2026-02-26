package messaging

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Envelope struct {
	ID            string          `json:"id"`
	CorrelationID string          `json:"correlationId"`
	Type          string          `json:"type"`
	Version       string          `json:"version"`
	Source        string          `json:"source"`
	CreatedAtUTC  time.Time       `json:"createdAtUtc"`
	Payload       json.RawMessage `json:"payload"`
}

func NewEnvelope(source string, correlationID string, messageType string, version string, payload []byte) Envelope {
	if correlationID == "" {
		correlationID = uuid.NewString()
	}
	if version == "" {
		version = "v1"
	}
	if len(payload) == 0 {
		payload = []byte("{}")
	}

	return Envelope{
		ID:            uuid.NewString(),
		CorrelationID: correlationID,
		Type:          messageType,
		Version:       version,
		Source:        source,
		CreatedAtUTC:  time.Now().UTC(),
		Payload:       payload,
	}
}

func (e Envelope) Validate() error {
	if e.ID == "" {
		return errors.New("envelope.id is required")
	}
	if e.CorrelationID == "" {
		return errors.New("envelope.correlationId is required")
	}
	if e.Type == "" {
		return errors.New("envelope.type is required")
	}
	if e.Source == "" {
		return errors.New("envelope.source is required")
	}
	if len(e.Payload) == 0 {
		return errors.New("envelope.payload is required")
	}
	return nil
}
