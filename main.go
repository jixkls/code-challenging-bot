package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type User struct {
    ID        int    `json:"id"`
    FirstName string `json:"first_name"`
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

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending reply:", err)
	}
	defer resp.Body.Close()
}

func webhookHandler(wri http.ResponseWriter, req *http.Request){
	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		return
	}
	defer req.Body.Close()

	var update Update
	err = json.Unmarshal(body, &update)
	if err != nil {
		fmt.Println("Error unmarshaling update:", err)
		return
	}

	userText := update.Message.Text
	chatID := update.Message.Chat.ID

	fmt.Printf("Received message: %s from chat ID: %d\n", userText, chatID)

	sendReply(chatID, "You said: "+userText)

}

func healthHandler (wri http.ResponseWriter, req *http.Request) {
	wri.WriteHeader(http.StatusOK)
	fmt.Fprint(wri, "Code Challenger Bot is running!")
}

func main() {
	godotenv.Load()

	http.HandleFunc("/", healthHandler)
	http.HandleFunc("/webhook", webhookHandler)

	fmt.Println("Starting server on :8080")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" //fallback para porta
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
