package config

import (
	"fmt"
	"log"
	"os"

	"ai-backend-go/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	LoadEnv()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	log.Println("✅ Connected to PostgreSQL successfully!")

	// Auto-migrate all models
	err = db.AutoMigrate(
		&models.Interview{},
		&models.User{},
		&models.Resume{},
		&models.Question{},
		&models.Answer{},
		&models.Feedback{},
	)
	if err != nil {
		log.Fatalf("❌ AutoMigration failed: %v", err)
	}

	DB = db
	log.Println("✅ Database migration completed.")
}
