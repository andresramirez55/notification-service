package database

import (
	"notification-service/config"
	"notification-service/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	cfg := config.LoadConfig()

	var err error

	// Determine database type from URL
	if isPostgreSQL(cfg.DatabaseURL) {
		DB, err = connectPostgreSQL(cfg.DatabaseURL)
	} else {
		// For production, use PostgreSQL
		if os.Getenv("PORT") != "" {
			log.Println("⚠️ Warning: Using SQLite in production is not recommended")
		}
		DB, err = connectSQLite(cfg.DatabaseURL)
	}

	if err != nil {
		return nil, err
	}

	// Auto migrate the schema
	err = DB.AutoMigrate(&models.Event{})
	if err != nil {
		log.Printf("Error migrating database: %v", err)
		return nil, err
	}

	log.Printf("✅ Database connected and migrated successfully using %s", getDBType(cfg.DatabaseURL))
	return DB, nil
}

func connectPostgreSQL(databaseURL string) (*gorm.DB, error) {
	log.Printf("Attempting to connect to PostgreSQL...")
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Silent to reduce noise
	})
	if err != nil {
		log.Printf("Failed to connect to PostgreSQL: %v", err)
		return nil, err
	}
	log.Println("✅ Successfully connected to PostgreSQL")
	return db, nil
}

func connectSQLite(databaseURL string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}

func isPostgreSQL(databaseURL string) bool {
	return len(databaseURL) > 11 && (databaseURL[:11] == "postgresql:" || databaseURL[:8] == "postgres:")
}

func getDBType(databaseURL string) string {
	if isPostgreSQL(databaseURL) {
		return "PostgreSQL"
	}
	return "SQLite"
}

func GetDB() *gorm.DB {
	return DB
}

