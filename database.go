package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db       *gorm.DB
	validate *validator.Validate
)

func InitDatabase() {
	requiredEnvVars := []string{
		"SUPABASE_HOST",
		"SUPABASE_USER",
		"SUPABASE_PASSWORD",
		"SUPABASE_DATABASE",
		"SUPABASE_PORT",
	}

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("Environment variable %s is required but not set", envVar)
		}
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=UTC",
		os.Getenv("SUPABASE_HOST"),
		os.Getenv("SUPABASE_USER"),
		os.Getenv("SUPABASE_PASSWORD"),
		os.Getenv("SUPABASE_DATABASE"),
		os.Getenv("SUPABASE_PORT"),
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully!")

	// Auto-migrate tables
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	validate = validator.New()
}

// Return User
func GetOrCreateUser(telegramID int64, username string) (*User, error) {
	var user User

	result := db.Where("telegram_id = ?", telegramID).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			user = User{
				TelegramID:       telegramID,
				Username:         username,
				Level:            1,
				ChallengesSolved: 0,
				CurrentStreak:    0,
				BestStreak:       0,
				CompletedToday:   false,
			}

			if err := db.Create(&user).Error; err != nil {
				return nil, fmt.Errorf("failed to create user: %v", err)
			}

			log.Printf("Created new user: %d (%s)", telegramID, username)  // FIX: Use telegramID
		} else {
			return nil, fmt.Errorf("database error: %v", result.Error)
		}
	}

	return &user, nil
}
