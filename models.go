package main

import "time"

// Database model for users (move from main.go and add GORM tags)
type UserStats struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	TelegramID       int64     `gorm:"uniqueIndex;not null" json:"telegram_id"`
	Username         string    `json:"username"`
	Level            int       `gorm:"default:1" json:"level"`
	ChallengesSolved int       `gorm:"default:0" json:"challenges_solved"`
	CurrentStreak    int       `gorm:"default:0" json:"current_streak"`
	BestStreak       int       `gorm:"default:0" json:"best_streak"`
	LastChallenge    time.Time `json:"last_challenge"`
	CompletedToday   bool      `gorm:"default:false" json:"completed_today"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Keep this struct (move from main.go, no changes needed)
type Challenge struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
	Example     string `json:"example"`
	Hint        string `json:"hint"`
}
