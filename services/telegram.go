package services

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// TelegramUser represents a Telegram user in incoming updates.
type TelegramUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

// Chat represents a Telegram chat.
type Chat struct {
	ID       int    `json:"id"`
	ChatType string `json:"type"`
}

// Message represents an incoming Telegram message.
type Message struct {
	MessageID int          `json:"message_id"`
	From      TelegramUser `json:"from"`
	Chat      Chat         `json:"chat"`
	Text      string       `json:"text"`
}

// Update represents a Telegram webhook update.
type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

// sendMessageBody is the JSON payload for Telegram's sendMessage API.
type sendMessageBody struct {
	ChatID    int    `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

// SendReply sends a text message to a Telegram chat.
// Uses InsecureSkipVerify as a workaround for container TLS issues on Fly.io.
func SendReply(chatID int, text string) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	url := "https://api.telegram.org/bot" + token + "/sendMessage"

	msgBody := sendMessageBody{
		ChatID: chatID,
		Text:   text,
	}

	jsonData, err := json.Marshal(msgBody)
	if err != nil {
		fmt.Println("Error marshaling reply:", err)
		return
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
