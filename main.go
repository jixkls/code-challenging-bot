package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type User struct {
    // Database fields (GORM)
    ID               uint      `gorm:"primaryKey" json:"db_id"`
    TelegramID       int64     `gorm:"uniqueIndex;not null" json:"id"`
    FirstName        string    `json:"first_name"`
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

type Chat struct {
    ID       int    `json:"id"`
    ChatType string `json:"type"`
}

type Message struct {
    MessageID int    `json:"message_id"`
    From      User   `json:"from"`
    Chat      Chat   `json:"chat"`
    Text      string `json:"text"`
}

//struct maior
type Update struct {
    UpdateID int     `json:"update_id"`
    Message  Message `json:"message"`
}

type SendMessage struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

func sendReply(ChatID int, Text string) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	url := "https://api.telegram.org/bot" + token + "/sendMessage"

	msgBody := SendMessage{
		ChatID: ChatID,
		Text:   Text,
	}

	jsonData, err := json.Marshal(msgBody)
	if err != nil {
		fmt.Println("Error marshaling reply:", err)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	res, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending reply:", err)
		return
	}
	defer res.Body.Close()
}

func webhookHandler(res http.ResponseWriter, req *http.Request){
	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		res.WriteHeader(http.StatusOK)
		return
	}
	defer req.Body.Close()
	fmt.Println(req.Body)

	if len(body) == 0 {
		fmt.Println("Empty request body")
		res.WriteHeader(http.StatusOK)
		return
	}

	var update Update
	err = json.Unmarshal(body, &update)
	if err != nil {
		fmt.Println("Error unmarshaling update:", err)
		res.WriteHeader(http.StatusOK)
		return
	}

	userText := update.Message.Text
	chatID := update.Message.Chat.ID

	isCommand := strings.HasPrefix(userText, "/")
	commands := []string{"/start", "/help", "/challenge", "/easy", "/medium", "/hard"}
	if isCommand {
		switch userText {
			case "/start":
					sendReply(chatID, "Bot iniciado!")
			case "/help":
					sendReply(chatID, "Comandos disponíveis: " + strings.Join(commands, ", "))
			case "/challenge":
					sendReply(chatID, "🎯 Random Challenge!\n\n**Easy Problem:**\nWrite a function that returns the sum of two integers.")
			case "/easy":
					sendReply(chatID, "🟢 Easy Challenge!\n\n**Problem:** Write a function that checks if a number is even.\n\n```go\nfunc isEven(n int) bool {\n    // Your code here\n}\n```")
			case "/medium":
					sendReply(chatID, "🟡 Medium Challenge!\n\n**Problem:** Write a function that reverses a string.\n\n```go\nfunc reverse(s string) string {\n    // Your code here\n}\n```")
			case "/hard":
					sendReply(chatID, "🔴 Hard Challenge!\n\n**Problem:** Implement a binary search function.\n\n```go\nfunc binarySearch(arr []int, target int) int {\n    // Your code here\n    // Return index if found, -1 if not found\n}\n```")
			default:
    			sendReply(chatID, "Unknown command. Type /help for available commands.")
}}
	res.WriteHeader(http.StatusOK)
}

func healthHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	fmt.Fprint(res, "Code Challenger Bot is running!")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println("Attempting to connect to database...")
	defer func() {
		if recovery := recover(); recovery != nil {
			log.Println("Recovered from panic:", recovery)
		}
	}() //Fechar isso se não fica reiniciando infinitamente

	InitDatabase()
	log.Println("Database connection established.")


	http.HandleFunc("/", healthHandler)
	http.HandleFunc("/webhook", webhookHandler)

	fmt.Println("Starting server on :8080")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" //fallback para porta
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
