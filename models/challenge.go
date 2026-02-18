package models

import "time"

// Challenge represents a Go programming challenge with multiple-choice answers.
type Challenge struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Difficulty  string     `json:"difficulty"`
	CodeSnippet string     `json:"code_snippet"`
	Options     [4]string  `json:"options"`
	CorrectIdx  int        `json:"correct_idx"` // 0-3 maps to A-D
	Hint        string     `json:"hint"`
}

// ActiveChallenge tracks a user's current in-progress challenge.
// Stored in memory — ephemeral by design (bot restart = new challenge).
type ActiveChallenge struct {
	UserTelegramID int64
	Challenge      *Challenge
	HintUsed       bool
	StartedAt      time.Time
}
