package notification

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/dbd-guard/backend/internal/config"
	"google.golang.org/api/option"
	"gopkg.in/gomail.v2"
)

type firebaseSMTPService struct {
	fcmClient *messaging.Client
	cfg       *config.AppConfig
}

// NewFirebaseSMTPService initializes the real notification service with FCM and Gomail
func NewFirebaseSMTPService(cfg *config.AppConfig) (Service, error) {
	ctx := context.Background()

	var opts []option.ClientOption
	if cfg.FirebaseCredentialsPath != "" {
		opts = append(opts, option.WithCredentialsFile(cfg.FirebaseCredentialsPath))
	}

	app, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		return nil, err
	}

	fcmClient, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	return &firebaseSMTPService{
		fcmClient: fcmClient,
		cfg:       cfg,
	}, nil
}

func (s *firebaseSMTPService) SendPushNotification(ctx context.Context, fcmToken, title, body string, data map[string]string) error {
	message := &messaging.Message{
		Token: fcmToken,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
	}

	var err error
	var response string
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		response, err = s.fcmClient.Send(ctx, message)
		if err == nil {
			log.Printf("[FCM] [AUDIT] SUCCESS: Successfully sent message to %s (Response: %s)\n", fcmToken, response)
			return nil
		}
		
		log.Printf("[FCM] [AUDIT] FAIL: Error sending message to %s (Attempt %d/%d): %v\n", fcmToken, i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(time.Duration((i+1)*2) * time.Second) // Simple backoff
		}
	}

	return err
}

func (s *firebaseSMTPService) SendEmail(ctx context.Context, toEmail, subject, body string) error {
	m := gomail.NewMessage()
	
	// If the config SMTPUser is empty, we fallback to log simulation
	if s.cfg.SMTPUser == "" || s.cfg.SMTPPass == "" {
		log.Printf("[SMTP SIMULATION] To: %s | Subject: %s | Body: %s", toEmail, subject, body)
		return nil
	}

	fromAddress := s.cfg.SMTPUser
	// If they use standard format, the user string might just be an email. If not, can wrap it if needed.
	if !strings.Contains(fromAddress, "@") {
		fromAddress = fromAddress + "@dbdguard.id" // Fallback dummy domain, better to be exact email
	}

	m.SetHeader("From", fromAddress)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	port, err := strconv.Atoi(s.cfg.SMTPPort)
	if err != nil {
		port = 587 // default
	}

	d := gomail.NewDialer(s.cfg.SMTPHost, port, s.cfg.SMTPUser, s.cfg.SMTPPass)
	
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		err = d.DialAndSend(m)
		if err == nil {
			log.Printf("[SMTP] [AUDIT] SUCCESS: Successfully sent email to %s\n", toEmail)
			return nil
		}
		
		log.Printf("[SMTP] [AUDIT] FAIL: Error sending email to %s (Attempt %d/%d): %v\n", toEmail, i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(time.Duration((i+1)*2) * time.Second) // Simple backoff
		}
	}

	return err
}
