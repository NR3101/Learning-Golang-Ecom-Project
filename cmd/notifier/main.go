package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/NR3101/go-ecom-project/internal/config"
	"github.com/NR3101/go-ecom-project/internal/models"
	"github.com/NR3101/go-ecom-project/internal/notications"
	"github.com/NR3101/go-ecom-project/internal/providers"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-aws/sqs"
	"github.com/ThreeDotsLabs/watermill/message"
)

func main() {
	fmt.Println("Starting Notification Service...")

	ctx := context.Background()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize email config
	emailCfg := &notications.SMTPConfig{
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		Username: cfg.SMTP.Username,
		Password: cfg.SMTP.Password,
		From:     cfg.SMTP.From,
	}

	emailSender := notications.NewEmailSender(emailCfg)

	// Create aws config for SQS
	awsCfg, err := providers.CreateAwsConfig(ctx, cfg.Aws.S3Endpoint, cfg.Aws.Region)
	if err != nil {
		log.Fatalf("Failed to create AWS config: %v", err)
	}

	// Create SQS subscriber
	logger := watermill.NewStdLogger(false, false)
	subscriber, err := sqs.NewSubscriber(sqs.SubscriberConfig{
		AWSConfig: awsCfg,
	}, logger)
	if err != nil {
		log.Fatalf("Failed to create SQS subscriber: %v", err)
	}

	// Subscribe to the SQS queue
	messages, err := subscriber.Subscribe(ctx, cfg.Aws.EventQueueName)
	if err != nil {
		subscriber.Close()
		log.Fatalf("Failed to subscribe to SQS queue: %v", err)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Notification Service started. Waiting for messages...")

	for {
		select {
		case msg := <-messages:
			if err := processMessage(msg, emailSender); err != nil {
				log.Printf("Failed to process message: %v", err)
				msg.Nack()
			} else {
				msg.Ack()
			}
		case sig := <-sigChan:
			log.Printf("Received signal: %v. Shutting down...", sig)
			subscriber.Close()
			return
		}
	}
}

func processMessage(msg *message.Message, emailSender *notications.EmailSender) interface{} {
	eventType := msg.Metadata.Get("event_type")
	log.Printf("Received message with event type: %s", eventType)

	switch eventType {
	case "user_authenticated":
		return handleUserAuthenticated(msg, emailSender)
	default:
		log.Printf("Unknown event type: %s", eventType)
		return nil
	}
}

func handleUserAuthenticated(msg *message.Message, emailSender *notications.EmailSender) error {
	var user models.User
	if err := json.Unmarshal(msg.Payload, &user); err != nil {
		return fmt.Errorf("failed to unmarshal message payload: %w", err)
	}

	userName := user.FirstName + " " + user.LastName
	if userName == "" {
		userName = user.Email
	}

	log.Println("Sending login notification email to:", user.Email)
	return emailSender.SendLoginNotification(user.Email, userName)
}
