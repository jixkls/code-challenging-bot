package services

import (
	"time"

	"telegram-bot/models"
)

// UpdateProgress updates a user's stats after solving a challenge.
// Returns true if the user leveled up.
func UpdateProgress(user *models.User, difficulty string) (bool, error) {
	user.ChallengesSolved++

	// Streak logic
	now := time.Now().UTC()
	today := now.Truncate(24 * time.Hour)

	if !user.LastChallenge.IsZero() {
		lastDay := user.LastChallenge.UTC().Truncate(24 * time.Hour)
		yesterday := today.AddDate(0, 0, -1)

		switch {
		case lastDay.Equal(today):
			// Already solved one today — no streak change
		case lastDay.Equal(yesterday):
			// Consecutive day — extend streak
			user.CurrentStreak++
		default:
			// Gap of 2+ days — reset streak
			user.CurrentStreak = 1
		}
	} else {
		// First ever challenge
		user.CurrentStreak = 1
	}

	if user.CurrentStreak > user.BestStreak {
		user.BestStreak = user.CurrentStreak
	}

	// Level progression: level N requires N*(N+1)/2 total challenges
	leveledUp := false
	nextLevelReq := user.Level * (user.Level + 1) / 2
	if user.ChallengesSolved >= nextLevelReq {
		user.Level++
		leveledUp = true
	}

	user.LastChallenge = now
	user.CompletedToday = true

	if err := SaveUser(user); err != nil {
		return false, err
	}

	return leveledUp, nil
}
