package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/NR3101/go-ecom-project/internal/providers"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-aws/sqs"
	"github.com/ThreeDotsLabs/watermill/message"

	appconfig "github.com/NR3101/go-ecom-project/internal/config"
	_ "github.com/aws/smithy-go/endpoints"
)

type EventPublisher struct {
	publisher message.Publisher
	queueName string
}

func NewEventPublisher(ctx context.Context, cfg *appconfig.AwsConfig) (*EventPublisher, error) {
	logger := watermill.NewStdLogger(false, false)

	awsConfig, err := providers.CreateAwsConfig(ctx, cfg.S3Endpoint, cfg.Region)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS config: %w", err)
	}

	publisherCfg := sqs.PublisherConfig{
		AWSConfig: awsConfig,
		Marshaler: nil,
	}

	publisher, err := sqs.NewPublisher(publisherCfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create SQS publisher: %w", err)
	}

	return &EventPublisher{
		publisher: publisher,
		queueName: cfg.EventQueueName,
	}, nil
}

func (ep *EventPublisher) Publish(eventType string, payload interface{}, metadata map[string]string) error {
	eventData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal event payload: %w", err)
	}

	msg := message.NewMessage(watermill.NewUUID(), eventData)

	msg.Metadata.Set("event_type", eventType)
	for key, value := range metadata {
		msg.Metadata.Set(key, value)
	}

	return ep.publisher.Publish(ep.queueName, msg)
}

func (ep *EventPublisher) Close() error {
	return ep.publisher.Close()
}
