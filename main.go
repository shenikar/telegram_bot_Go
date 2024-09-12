package main

import (
	"log"
	"telegram_bot_go/adapter"
	"telegram_bot_go/config"
	"telegram_bot_go/repository"
	"telegram_bot_go/service"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := repository.GetConnect(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to the DB")
	}

	hashService := service.NewHashService()
	userRepo := repository.NewUserRepo(db)
	userService := service.NewUserService(userRepo, cfg)

	botApi, err := tgbot.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Fatal(err)
	}

	botApi.Debug = true
	log.Printf("Authorized on account %s", botApi.Self.UserName)

	bot := adapter.NewBot(botApi, hashService, userService, cfg)
	bot.Start()
}
