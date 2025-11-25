package handlers

import (
	"net/http"
	"notification-service/services"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	emailService *services.EmailService
}

func NewNotificationHandler(emailService *services.EmailService) *NotificationHandler {
	return &NotificationHandler{
		emailService: emailService,
	}
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
