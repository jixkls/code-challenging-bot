package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"telegram-bot/handlers"
	"telegram-bot/services"

	"github.com/joho/godotenv"
)

func healthHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	fmt.Fprint(res, "Code Challenger Bot is running!")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	defer func() {
		if recovery := recover(); recovery != nil {
			log.Println("Recovered from panic:", recovery)
		}
	}()

	log.Println("Attempting to connect to database...")
	services.InitDatabase()
	log.Println("Database connection established.")

	if os.Getenv("GEMINI_API_KEY") == "" {
		log.Println("Warning: GEMINI_API_KEY not set. Using static challenges only.")
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/webhook", handlers.WebhookHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Starting server on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
