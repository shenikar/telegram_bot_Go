package main

import (
	"log"
	"os"
	"telegram_bot_go/adapter"
	"telegram_bot_go/repository"
	"telegram_bot_go/service"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file")
	}

	db, err := repository.GetConnect()
	if err != nil {
		log.Fatal("Failed to connect to the DB")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("Error: TELEGRAM_BOT_TOKEN is not set")
	}

	botApi, err := tgbot.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	botApi.Debug = true
	log.Printf("Authorized on account %s", botApi.Self.UserName)

	hashService := service.NewHashService()
	userRepo := repository.NewUserRepo(db)

	bot := adapter.NewBot(botApi, hashService, userRepo)
	bot.Start()
}
