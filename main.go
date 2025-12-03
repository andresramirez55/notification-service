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

	// Setup router FIRST so health check works even if DB fails
	router := gin.Default()

	// Health check endpoint - must work even if DB is down
	router.GET("/health", func(c *gin.Context) {
		// Try to check DB connection
		dbStatus := "unknown"
		if db, err := database.GetDB(); err == nil && db != nil {
			sqlDB, err := db.DB()
			if err == nil {
				if err := sqlDB.Ping(); err == nil {
					dbStatus = "connected"
				} else {
					dbStatus = "disconnected"
				}
			}
		} else {
			dbStatus = "not_initialized"
		}

		statusCode := 200
		if dbStatus != "connected" {
			statusCode = 503 // Service Unavailable
		}

		c.JSON(statusCode, gin.H{
			"status":      "ok",
			"service":     "notification-service",
			"version":     "1.0.0",
			"database":    dbStatus,
			"time":        time.Now().Format(time.RFC3339),
		})
	})

	// Initialize email service (needed for API endpoints and scheduler)
	emailService := services.NewEmailService(cfg)

	// Initialize handlers
	notificationHandler := handlers.NewNotificationHandler(emailService)

	// Initialize database
	log.Println("🔌 Connecting to database...")
	db, err := database.InitDB()
	var schedulerService *services.SchedulerService
	if err != nil {
		log.Printf("❌ Failed to connect to database: %v", err)
		log.Println("⚠️ Service will start but scheduler will not work until database is available")
		log.Println("⚠️ Health check will return 503 until database connection is established")
		// Don't fatal - let the service start so health check works
	} else {
		// Initialize repository
		eventRepo := repositories.NewEventRepository(db)

		// Initialize scheduler service
		log.Println("⏰ Initializing notification scheduler...")
		schedulerService = services.NewSchedulerService(eventRepo, emailService)
		schedulerService.Start()
		log.Println("✅ Notification scheduler started")

		// Link scheduler to handler for manual trigger endpoint
		notificationHandler.SetScheduler(schedulerService)
		log.Println("✅ Notification scheduler linked to handler")
	}

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
		api.POST("/notifications/check", notificationHandler.CheckNotificationsNow)
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
	if db != nil {
		log.Println("✅ Scheduler is running - checking for events every hour")
	} else {
		log.Println("⚠️ Scheduler is NOT running - database connection required")
	}
	log.Println("🔍 Test with: curl http://localhost:" + port + "/health")
	log.Println("========================================")

	if err := router.Run(":" + port); err != nil {
		log.Fatal("❌ Failed to start notification service:", err)
	}
}

