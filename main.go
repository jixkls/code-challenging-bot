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

	http.HandleFunc("/", healthHandler)
	http.HandleFunc("/webhook", handlers.WebhookHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Starting server on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
