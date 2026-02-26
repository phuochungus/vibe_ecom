package messaging

import (
	"context"
	"log"
)

type Consumer struct {
	serviceName  string
	kafkaBrokers string
	rabbitMQURL  string
}

func NewConsumer(serviceName string, kafkaBrokers string, rabbitMQURL string) *Consumer {
	return &Consumer{
		serviceName:  serviceName,
		kafkaBrokers: kafkaBrokers,
		rabbitMQURL:  rabbitMQURL,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Printf("[%s] consumer started (kafka=%s rabbitmq=%s)", c.serviceName, c.kafkaBrokers, c.rabbitMQURL)
	<-ctx.Done()
	log.Printf("[%s] consumer stopped", c.serviceName)
}
