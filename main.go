package main

import (
	"log"
	"os"
	"time"

	"notification-service/config"
	"notification-service/database"
	"notification-service/handlers"
	"notification-service/repositories"
	"notification-service/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("========================================")
	log.Println("📧 NOTIFICATION SERVICE - STARTING...")
	log.Println("========================================")
	log.Println("⏰ Timestamp:", time.Now().Format(time.RFC3339))

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	log.Println("🔌 Connecting to database...")
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Initialize repository
	eventRepo := repositories.NewEventRepository(db)

	// Initialize email service
	emailService := services.NewEmailService(cfg)

	// Initialize scheduler service
	log.Println("⏰ Initializing notification scheduler...")
	schedulerService := services.NewSchedulerService(eventRepo, emailService)
	schedulerService.Start()
	log.Println("✅ Notification scheduler started")

	// Initialize handlers
	notificationHandler := handlers.NewNotificationHandler(emailService)

	// Setup router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "notification-service",
			"version": "1.0.0",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Configure CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(corsConfig))

	// API routes
	api := router.Group("/api/v1")
	{
		api.POST("/notifications/email", notificationHandler.SendEmail)
		api.GET("/notifications/status", notificationHandler.GetStatus)
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Println("========================================")
	log.Printf("📧 Notification service starting on port %s", port)
	log.Println("========================================")
	log.Println("✅ Ready to accept notification requests!")
	log.Println("✅ Scheduler is running - checking for events every hour")
	log.Println("🔍 Test with: curl http://localhost:" + port + "/health")
	log.Println("========================================")

	if err := router.Run(":" + port); err != nil {
		log.Fatal("❌ Failed to start notification service:", err)
	}
}

