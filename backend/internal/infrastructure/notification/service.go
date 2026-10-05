package notification

import (
	"context"
	"log"

	"github.com/dbd-guard/backend/internal/config"
)

// Service defines the interface for sending notifications
type Service interface {
	SendPushNotification(ctx context.Context, fcmToken, title, body string, data map[string]string) error
	SendEmail(ctx context.Context, toEmail, subject, body string) error
}

type mockNotificationService struct{}

// NewNotificationService returns the real notification service if configured,
// otherwise it falls back to the mock implementation for development
func NewNotificationService(cfg *config.AppConfig) Service {
	// Only instantiate the real service if Firebase credentials are provided
	// Otherwise, use mock to prevent application crash during local development
	if cfg.FirebaseCredentialsPath != "" {
		svc, err := NewFirebaseSMTPService(cfg)
		if err == nil {
			log.Println("Initialized Real Firebase/SMTP Notification Service")
			return svc
		}
		log.Printf("Failed to initialize Firebase/SMTP Service, falling back to mock: %v\n", err)
	}

	log.Println("Initialized Mock Notification Service")
	return &mockNotificationService{}
}

func (s *mockNotificationService) SendPushNotification(ctx context.Context, fcmToken, title, body string, data map[string]string) error {
	log.Printf("[MOCK FCM] Sending push to %s | Title: %s | Body: %s\n", fcmToken, title, body)
	return nil
}

func (s *mockNotificationService) SendEmail(ctx context.Context, toEmail, subject, body string) error {
	log.Printf("[MOCK SMTP] Sending email to %s | Subject: %s\n", toEmail, subject)
	return nil
}
