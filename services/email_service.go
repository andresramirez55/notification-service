package services

import (
	"notification-service/config"
	"fmt"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailService struct {
	cfg *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	// Log configuration status
	log.Println("📧 Email Service Configuration:")
	if cfg.SendGridAPIKey != "" {
		log.Printf("  ✅ SendGrid API Key: Configured (from email: %s)", cfg.FromEmail)
	} else {
		log.Println("  ⚠️ SendGrid API Key: NOT configured - Email notifications will fail")
		log.Println("  💡 To enable email notifications, set SENDGRID_API_KEY environment variable")
	}

	return &EmailService{
		cfg: cfg,
	}
}

type EmailRequest struct {
	To          string `json:"to" binding:"required"`
	Subject     string `json:"subject" binding:"required"`
	Body        string `json:"body" binding:"required"`
	RecipientName string `json:"recipient_name,omitempty"`
}

func (s *EmailService) SendEmail(req *EmailRequest) error {
	if s.cfg.SendGridAPIKey == "" {
		return fmt.Errorf("SendGrid API key not configured")
	}

	from := mail.NewEmail("Calendar Reminder", s.cfg.FromEmail)
	
	recipientName := req.RecipientName
	if recipientName == "" {
		recipientName = "User"
	}
	to := mail.NewEmail(recipientName, req.To)

	message := mail.NewSingleEmail(from, req.Subject, to, req.Body, req.Body)
	client := sendgrid.NewSendClient(s.cfg.SendGridAPIKey)

	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	log.Printf("✅ Email sent successfully to %s, status: %d", req.To, response.StatusCode)
	return nil
}

