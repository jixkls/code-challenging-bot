package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"telegram-bot/models"
	"telegram-bot/services"
)

// WebhookHandler processes incoming Telegram webhook updates.
func WebhookHandler(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		res.WriteHeader(http.StatusOK)
		return
	}
	defer req.Body.Close()

	if len(body) == 0 {
		fmt.Println("Empty request body")
		res.WriteHeader(http.StatusOK)
		return
	}

	var update services.Update
	err = json.Unmarshal(body, &update)
	if err != nil {
		fmt.Println("Error unmarshaling update:", err)
		res.WriteHeader(http.StatusOK)
		return
	}

	userText := strings.TrimSpace(update.Message.Text)
	chatID := update.Message.Chat.ID
	telegramID := update.Message.From.ID

	// Get or create user in database
	user, err := services.GetOrCreateUser(telegramID, update.Message.From.Username)
	if err != nil {
		log.Printf("Database error: %v", err)
		services.SendReply(chatID, "Welcome! (Database temporarily unavailable)")
		res.WriteHeader(http.StatusOK)
		return
	}

	// Check if this is an answer (A, B, C, or D)
	upperText := strings.ToUpper(userText)
	if len(upperText) == 1 && upperText >= "A" && upperText <= "D" {
		handleAnswer(chatID, telegramID, upperText, user)
		res.WriteHeader(http.StatusOK)
		return
	}

	// Command routing
	if strings.HasPrefix(userText, "/") {
		command := strings.Split(userText, "@")[0] // Handle /command@botname
		switch command {
		case "/start":
			handleStart(chatID, user)
		case "/help":
			handleHelp(chatID)
		case "/challenge":
			handleChallenge(chatID, telegramID, "")
		case "/easy":
			handleChallenge(chatID, telegramID, "easy")
		case "/medium":
			handleChallenge(chatID, telegramID, "medium")
		case "/hard":
			handleChallenge(chatID, telegramID, "hard")
		case "/hint":
			handleHint(chatID, telegramID)
		case "/skip":
			handleSkip(chatID, telegramID)
		case "/stats":
			handleStats(chatID, user)
		default:
			services.SendReply(chatID, "Unknown command. Type /help for available commands.")
		}
	}

	res.WriteHeader(http.StatusOK)
}

func handleStart(chatID int, user *models.User) {
	message := fmt.Sprintf(
		"Welcome %s!\n\n"+
			"I'm Code Challenger Bot — solve Go challenges, build streaks, and level up!\n\n"+
			"Your Progress:\n"+
			"  Level: %d\n"+
			"  Challenges Solved: %d\n"+
			"  Current Streak: %d days\n\n"+
			"Type /easy, /medium, or /hard to get a challenge!",
		user.Username, user.Level, user.ChallengesSolved, user.CurrentStreak,
	)
	services.SendReply(chatID, message)
}

func handleHelp(chatID int) {
	help := "Available commands:\n\n" +
		"/easy — Get an easy Go challenge\n" +
		"/medium — Get a medium Go challenge\n" +
		"/hard — Get a hard Go challenge\n" +
		"/challenge — Get a random difficulty challenge\n" +
		"/hint — Get a hint for your current challenge\n" +
		"/skip — Skip the current challenge\n" +
		"/stats — View your detailed progress\n" +
		"/help — Show this message\n\n" +
		"Answer with A, B, C, or D after receiving a challenge!"
	services.SendReply(chatID, help)
}

func handleChallenge(chatID int, telegramID int64, difficulty string) {
	// Check if user already has an active challenge
	if active := services.GetActiveChallenge(telegramID); active != nil {
		services.SendReply(chatID,
			"You already have an active challenge! Answer with A/B/C/D, use /hint, or /skip it.")
		return
	}

	// Pick random difficulty if none specified
	if difficulty == "" {
		difficulties := []string{"easy", "medium", "hard"}
		difficulty = difficulties[rand.Intn(len(difficulties))]
	}

	challenge := services.GetRandomChallenge(difficulty)
	if challenge == nil {
		services.SendReply(chatID, "No challenges available for that difficulty. Try another!")
		return
	}

	services.SetActiveChallenge(telegramID, challenge)

	difficultyLabel := map[string]string{
		"easy":   "[EASY]",
		"medium": "[MEDIUM]",
		"hard":   "[HARD]",
	}

	message := fmt.Sprintf(
		"%s %s\n\n%s\n\n%s\n\n%s\n%s\n%s\n%s\n\nReply with A, B, C, or D!",
		difficultyLabel[challenge.Difficulty],
		challenge.Title,
		challenge.Description,
		challenge.CodeSnippet,
		challenge.Options[0],
		challenge.Options[1],
		challenge.Options[2],
		challenge.Options[3],
	)
	services.SendReply(chatID, message)
}

func handleHint(chatID int, telegramID int64) {
	active := services.GetActiveChallenge(telegramID)
	if active == nil {
		services.SendReply(chatID, "No active challenge. Use /easy, /medium, or /hard to get one!")
		return
	}

	active.HintUsed = true
	services.SendReply(chatID, "Hint: "+active.Challenge.Hint)
}

func handleSkip(chatID int, telegramID int64) {
	active := services.GetActiveChallenge(telegramID)
	if active == nil {
		services.SendReply(chatID, "No active challenge to skip.")
		return
	}

	correctLetter := string(rune('A' + active.Challenge.CorrectIdx))
	services.ClearActiveChallenge(telegramID)
	services.SendReply(chatID, fmt.Sprintf(
		"Challenge skipped. The correct answer was %s.\n\nUse /easy, /medium, or /hard for a new challenge!",
		correctLetter,
	))
}

func handleStats(chatID int, user *models.User) {
	// Build XP progress bar
	currentLevel := user.Level
	currentReq := currentLevel * (currentLevel + 1) / 2
	prevReq := (currentLevel - 1) * currentLevel / 2
	progress := user.ChallengesSolved - prevReq
	needed := currentReq - prevReq

	barLen := 10
	filled := 0
	if needed > 0 {
		filled = progress * barLen / needed
	}
	if filled > barLen {
		filled = barLen
	}
	bar := strings.Repeat("#", filled) + strings.Repeat("-", barLen-filled)

	message := fmt.Sprintf(
		"Your Stats\n\n"+
			"Level: %d\n"+
			"Progress: [%s] %d/%d\n"+
			"Total Solved: %d\n"+
			"Current Streak: %d days\n"+
			"Best Streak: %d days\n\n"+
			"Keep going! Type /easy, /medium, or /hard for your next challenge.",
		user.Level, bar, progress, needed,
		user.ChallengesSolved,
		user.CurrentStreak,
		user.BestStreak,
	)
	services.SendReply(chatID, message)
}

func handleAnswer(chatID int, telegramID int64, answer string, user *models.User) {
	active := services.GetActiveChallenge(telegramID)
	if active == nil {
		services.SendReply(chatID, "No active challenge. Use /easy, /medium, or /hard to get one!")
		return
	}

	answerIdx := int(answer[0] - 'A') // A=0, B=1, C=2, D=3

	if answerIdx == active.Challenge.CorrectIdx {
		// Correct answer
		services.ClearActiveChallenge(telegramID)

		leveledUp, err := services.UpdateProgress(user, active.Challenge.Difficulty)
		if err != nil {
			log.Printf("Error updating progress: %v", err)
			services.SendReply(chatID, "Correct! But there was an error saving your progress.")
			return
		}

		message := fmt.Sprintf("Correct! Well done!\n\nChallenges Solved: %d\nStreak: %d days",
			user.ChallengesSolved, user.CurrentStreak)

		if leveledUp {
			message += fmt.Sprintf("\n\nLEVEL UP! You're now level %d!", user.Level)
		}

		message += "\n\nUse /easy, /medium, or /hard for another challenge!"
		services.SendReply(chatID, message)
	} else {
		// Wrong answer
		services.SendReply(chatID,
			"Not quite! Try again, use /hint for a clue, or /skip to move on.")
	}
}
