package services

import (
	"fmt"
	"log"
	"os"
	"time"

	"telegram-bot/models"

	"github.com/go-playground/validator/v10"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db       *gorm.DB
	validate *validator.Validate
)

// InitDatabase connects to Supabase (PostgreSQL) and runs auto-migration.
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

	if os.Getenv("RUN_MIGRATIONS") == "true" {
		log.Println("Running database migrations...")
		err = db.AutoMigrate(&models.User{})
		if err != nil {
			log.Fatal("Failed to migrate database:", err)
		}
		log.Println("Migrations complete.")
	} else {
		log.Println("Skipping migrations. Set RUN_MIGRATIONS=true to run.")
	}

	validate = validator.New()
}

// GetOrCreateUser retrieves an existing user or creates a new one.
// Also performs a lazy daily reset: if LastChallenge is from a previous day,
// CompletedToday is reset to false.
func GetOrCreateUser(telegramID int64, username string) (*models.User, error) {
	var user models.User

	result := db.Where("telegram_id = ?", telegramID).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			user = models.User{
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

			log.Printf("Created new user: %d (%s)", telegramID, username)
			return &user, nil
		}
		return nil, fmt.Errorf("database error: %v", result.Error)
	}

	// Lazy daily reset: if last challenge was on a different day, reset CompletedToday
	if !user.LastChallenge.IsZero() {
		today := time.Now().UTC().Truncate(24 * time.Hour)
		lastDay := user.LastChallenge.UTC().Truncate(24 * time.Hour)
		if !today.Equal(lastDay) {
			user.CompletedToday = false
			db.Save(&user)
		}
	}

	return &user, nil
}

// SaveUser persists updated user data to the database.
func SaveUser(user *models.User) error {
	if err := db.Save(user).Error; err != nil {
		return fmt.Errorf("failed to save user: %v", err)
	}
	return nil
}
