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

	"github.com/joho/godotenv"
)

type TelegramUser struct {
		ID               int64     `json:"id"`
    FirstName        string    `json:"first_name"`
    Username         string    `json:"username"`
}

type Chat struct {
    ID       int    `json:"id"`
    ChatType string `json:"type"`
}

type Message struct {
    MessageID int    `json:"message_id"`
    From      TelegramUser   `json:"from"`
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

	// Get or create user in database
	user, err := GetOrCreateUser(update.Message.From.ID, update.Message.From.Username)
	if err != nil {
		log.Printf("Database error: %v", err)
		sendReply(chatID, "Welcome! (Database temporarily unavailable)")
		return
		}

	isCommand := strings.HasPrefix(userText, "/")
	commands := []string{"/start", "/help", "/challenge", "/easy", "/medium", "/hard"}
	if isCommand {
		switch userText {
			case "/start":
    		// Send personalized welcome message
    		message := fmt.Sprintf("Welcome %s! 🎯\n\n📊 Your Progress:\n🔥 Level: %d\n✅ Challenges Solved: %d\n⚡ Current Streak: %d",
        user.Username, user.Level, user.ChallengesSolved, user.CurrentStreak)
					sendReply(chatID, message)

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
		log.Println("Error loading .env file")
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
