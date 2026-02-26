package messaging

import (
	"context"
	"encoding/json"
	"log"

	platformmsg "golf-store/be/internal/platform/messaging"
)

type CommandPublisher interface {
	PublishCommand(ctx context.Context, correlationID string, messageType string, payload any) (platformmsg.Envelope, error)
}

type BrokerPublisher struct {
	serviceName string
}

func NewBrokerPublisher(serviceName string) *BrokerPublisher {
	return &BrokerPublisher{serviceName: serviceName}
}

func (p *BrokerPublisher) PublishCommand(ctx context.Context, correlationID string, messageType string, payload any) (platformmsg.Envelope, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return platformmsg.Envelope{}, err
	}

	envelope := platformmsg.NewEnvelope(p.serviceName, correlationID, messageType, "v1", body)
	log.Printf("[%s] publish command: type=%s correlationId=%s payload=%s", p.serviceName, envelope.Type, envelope.CorrelationID, string(envelope.Payload))

	_ = ctx
	// Skeleton only: producer integration to RabbitMQ/Kafka will be implemented in next phase.
	return envelope, nil
}
