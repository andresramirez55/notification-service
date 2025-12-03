package handlers

import (
	"net/http"
	"notification-service/services"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	emailService *services.EmailService
	scheduler    *services.SchedulerService
}

func NewNotificationHandler(emailService *services.EmailService) *NotificationHandler {
	return &NotificationHandler{
		emailService: emailService,
		scheduler:    nil, // Se setea después si está disponible
	}
}

func (h *NotificationHandler) SetScheduler(scheduler *services.SchedulerService) {
	h.scheduler = scheduler
}

func (h *NotificationHandler) SendEmail(c *gin.Context) {
	var req services.EmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	if err := h.emailService.SendEmail(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send email",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email sent successfully",
		"status":  "success",
		"to":      req.To,
	})
}

func (h *NotificationHandler) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service": "notification-service",
		"status":  "active",
		"features": gin.H{
			"email_notifications": true,
		},
	})
}

// CheckNotificationsNow fuerza la verificación de notificaciones manualmente
func (h *NotificationHandler) CheckNotificationsNow(c *gin.Context) {
	if h.scheduler == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "Scheduler not available",
			"message": "Database connection required for scheduler to work",
		})
		return
	}

	// Llamar al método del scheduler para verificar ahora
	h.scheduler.CheckAndSendNotifications()

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification check triggered successfully",
		"status":  "success",
	})
}
