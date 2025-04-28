package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"xm-exercise/pkg/models"
)

const (
	TopicCompanyCreated = "company.created"
	TopicCompanyUpdated = "company.updated"
	TopicCompanyDeleted = "company.deleted"
)

// Event represents a Kafka event
type Event struct {
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// KafkaProducer handles publishing events to Kafka
type KafkaProducer struct {
	writer *kafka.Writer
}

// NewKafkaProducer creates a new Kafka producer
func NewKafkaProducer(brokers []string) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaProducer{writer: writer}
}

// PublishCompanyCreated publishes a company created event
func (p *KafkaProducer) PublishCompanyCreated(company models.Company) error {
	return p.publishEvent(TopicCompanyCreated, "company.created", company)
}

// PublishCompanyUpdated publishes a company updated event
func (p *KafkaProducer) PublishCompanyUpdated(company *models.Company) error {
	return p.publishEvent(TopicCompanyUpdated, "company.updated", company)
}

// PublishCompanyDeleted publishes a company deleted event
func (p *KafkaProducer) PublishCompanyDeleted(companyID string) error {
	return p.publishEvent(TopicCompanyDeleted, "company.deleted", map[string]string{"id": companyID})
}

// publishEvent publishes an event to Kafka
func (p *KafkaProducer) publishEvent(topic, eventType string, data interface{}) error {
	event := Event{
		Type:      eventType,
		Timestamp: time.Now(),
		Data:      data,
	}

	value, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error marshaling event: %w", err)
	}

	err = p.writer.WriteMessages(context.Background(),
		kafka.Message{
			Topic: topic,
			Value: value,
		},
	)

	if err != nil {
		return fmt.Errorf("error publishing event to Kafka: %w", err)
	}

	return nil
}

// Close closes the Kafka writer
func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
