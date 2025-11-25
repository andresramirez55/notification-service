package services

import (
	"encoding/json"
	"fmt"
	"log"
	"notification-service/models"
	"notification-service/repositories"
	"time"

	"github.com/robfig/cron/v3"
)

type SchedulerService struct {
	eventRepo    *repositories.EventRepository
	emailService *EmailService
	cron         *cron.Cron
}

type FamilyMember struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Role  string `json:"role"`
}

func NewSchedulerService(eventRepo *repositories.EventRepository, emailService *EmailService) *SchedulerService {
	return &SchedulerService{
		eventRepo:    eventRepo,
		emailService: emailService,
		cron:         cron.New(cron.WithLocation(time.UTC)),
	}
}

// Start inicia el scheduler
func (s *SchedulerService) Start() {
	log.Println("🚀 Starting notification scheduler...")

	// Ejecutar inmediatamente al inicio
	go s.checkAndSendNotifications()

	// Programar verificación cada hora
	s.cron.AddFunc("0 * * * *", s.checkAndSendNotifications)

	// También programar verificación diaria a las 9:00 AM UTC para eventos del día siguiente
	s.cron.AddFunc("0 9 * * *", s.checkDayBeforeReminders)

	s.cron.Start()
	log.Println("✅ Notification scheduler started - checking every hour")
}

// Stop detiene el scheduler
func (s *SchedulerService) Stop() {
	if s.cron != nil {
		s.cron.Stop()
		log.Println("🛑 Notification scheduler stopped")
	}
}

// checkAndSendNotifications verifica y envía notificaciones
func (s *SchedulerService) checkAndSendNotifications() {
	log.Println("🔍 Checking for notifications to send...")

	// Verificar eventos para mañana (día anterior)
	s.checkDayBeforeReminders()

	// Verificar eventos para hoy (mismo día)
	s.checkSameDayReminders()
}

// checkDayBeforeReminders verifica eventos para mañana
func (s *SchedulerService) checkDayBeforeReminders() {
	log.Println("📅 Checking for day-before reminders (tomorrow's events)...")

	events, err := s.eventRepo.GetEventsForTomorrow()
	if err != nil {
		log.Printf("❌ Error getting events for tomorrow: %v", err)
		return
	}

	if len(events) == 0 {
		log.Println("ℹ️ No events found for tomorrow")
		return
	}

	log.Printf("📧 Found %d events for tomorrow, sending notifications...", len(events))

	for _, event := range events {
		s.sendEventNotification(event, "day_before")
	}
}

// checkSameDayReminders verifica eventos para hoy
func (s *SchedulerService) checkSameDayReminders() {
	log.Println("📅 Checking for same-day reminders (today's events)...")

	events, err := s.eventRepo.GetEventsForToday()
	if err != nil {
		log.Printf("❌ Error getting events for today: %v", err)
		return
	}

	if len(events) == 0 {
		log.Println("ℹ️ No events found for today")
		return
	}

	log.Printf("📧 Found %d events for today, sending notifications...", len(events))

	for _, event := range events {
		// Verificar si es hora de enviar (1 hora antes del evento)
		if s.shouldSendReminderNow(event) {
			s.sendEventNotification(event, "same_day")
		}
	}
}

// shouldSendReminderNow verifica si es hora de enviar el recordatorio (1 hora antes)
func (s *SchedulerService) shouldSendReminderNow(event *models.Event) bool {
	if event.Time == "" || event.IsAllDay {
		// Para eventos de todo el día, enviar en la mañana
		now := time.Now()
		return now.Hour() >= 8 && now.Hour() < 9
	}

	// Parsear hora del evento
	eventTime, err := time.Parse("15:04", event.Time)
	if err != nil {
		log.Printf("⚠️ Error parsing time for event %d: %v", event.ID, err)
		return true // En caso de error, enviar de todas formas
	}

	// Crear fecha/hora del evento para hoy
	today := time.Now()
	eventDateTime := time.Date(today.Year(), today.Month(), today.Day(),
		eventTime.Hour(), eventTime.Minute(), 0, 0, today.Location())

	// Enviar 1 hora antes del evento
	reminderTime := eventDateTime.Add(-1 * time.Hour)
	now := time.Now()

	// Verificar si estamos en la ventana de tiempo (dentro de la última hora)
	return now.After(reminderTime) && now.Before(eventDateTime)
}

// sendEventNotification envía notificaciones para un evento
func (s *SchedulerService) sendEventNotification(event *models.Event, reminderType string) {
	log.Printf("📧 Sending %s notification for event: %s (ID: %d)", reminderType, event.Title, event.ID)

	// Enviar email al usuario principal
	if event.Email != "" {
		subject, body := s.buildEmailContent(event, reminderType, "")
		if err := s.emailService.SendEmail(&EmailRequest{
			To:            event.Email,
			Subject:       subject,
			Body:          body,
			RecipientName: "",
		}); err != nil {
			log.Printf("❌ Error sending email to %s: %v", event.Email, err)
		} else {
			log.Printf("✅ Email sent to %s for event: %s", event.Email, event.Title)
		}
	}

	// Enviar notificaciones familiares si está habilitado
	if event.NotifyFamily {
		s.sendFamilyNotifications(event, reminderType)
	}
}

// sendFamilyNotifications envía notificaciones a miembros de la familia
func (s *SchedulerService) sendFamilyNotifications(event *models.Event, reminderType string) {
	var familyMembers []FamilyMember
	if event.FamilyMembers != "" {
		if err := json.Unmarshal([]byte(event.FamilyMembers), &familyMembers); err != nil {
			log.Printf("⚠️ Error parsing family members: %v", err)
			return
		}
	}

	var selectedChildren []string
	if event.SelectedChildren != "" {
		if err := json.Unmarshal([]byte(event.SelectedChildren), &selectedChildren); err != nil {
			log.Printf("⚠️ Error parsing selected children: %v", err)
		}
	}

	var recipients []FamilyMember

	// Filtrar por rol
	if event.NotifyPapa {
		for _, member := range familyMembers {
			if member.Role == "papa" {
				recipients = append(recipients, member)
			}
		}
	}

	if event.NotifyMama {
		for _, member := range familyMembers {
			if member.Role == "mama" {
				recipients = append(recipients, member)
			}
		}
	}

	// Enviar emails a los destinatarios
	for _, recipient := range recipients {
		if recipient.Email != "" {
			subject, body := s.buildEmailContent(event, reminderType, recipient.Name)
			childrenInfo := ""
			if len(selectedChildren) > 0 {
				childrenInfo = fmt.Sprintf("\n\nEste evento es para: %v", selectedChildren)
			}
			body += childrenInfo

			if err := s.emailService.SendEmail(&EmailRequest{
				To:            recipient.Email,
				Subject:       subject,
				Body:          body,
				RecipientName: recipient.Name,
			}); err != nil {
				log.Printf("❌ Error sending family email to %s: %v", recipient.Email, err)
			} else {
				log.Printf("✅ Family email sent to %s (%s)", recipient.Email, recipient.Name)
			}
		}
	}
}

// buildEmailContent construye el contenido del email
func (s *SchedulerService) buildEmailContent(event *models.Event, reminderType, recipientName string) (string, string) {
	var subject string
	var greeting string

	if recipientName != "" {
		greeting = fmt.Sprintf("Hola %s!", recipientName)
	} else {
		greeting = "Hola!"
	}

	switch reminderType {
	case "day_before":
		subject = fmt.Sprintf("Recordatorio: %s mañana", event.Title)
		body := fmt.Sprintf(`
%s

Te recordamos que mañana tenés:

Evento: %s
Fecha: %s
Hora: %s
%s

¡No te lo pierdas!
		`, greeting, event.Title, event.Date.Format("02/01/2006"), event.Time,
			func() string {
				if event.Location != "" {
					return fmt.Sprintf("Ubicación: %s", event.Location)
				}
				return ""
			}())
		return subject, body

	case "same_day":
		subject = fmt.Sprintf("Recordatorio: %s hoy", event.Title)
		body := fmt.Sprintf(`
%s

Te recordamos que hoy tenés:

Evento: %s
Hora: %s
%s

¡Que tengas un buen día!
		`, greeting, event.Title, event.Time,
			func() string {
				if event.Location != "" {
					return fmt.Sprintf("Ubicación: %s", event.Location)
				}
				return ""
			}())
		return subject, body

	default:
		return "Recordatorio de evento", "Tienes un evento programado"
	}
}
