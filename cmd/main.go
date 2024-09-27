package main

import (
	"log"
	"net/http"
	"telegram_bot_go/adapter/database"
	server "telegram_bot_go/adapter/http_server"
	rabbitmq "telegram_bot_go/adapter/queue"
	"telegram_bot_go/adapter/telegram"
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

	db, err := database.GetConnect(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to the DB")
	}
	rabbitMQService, err := rabbitmq.NewRabbitMqService(cfg.RabbitMQ.URL, cfg.RabbitMQ.Queue)
	if err != nil {
		log.Fatal("Failed to initialize RabbitMQ service")
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

	statsService := service.NewStatsService(userRepo)

	go func() {
		httpHandler := server.NewStatsHandler(statsService)
		http.Handle("/stats", httpHandler)
		log.Println("Starting HTTP server on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal("HTTP server failed:", err)
		}
	}()

	bot := telegram.NewBot(botApi, hashService, userService, rabbitMQService, statsService, cfg)
	bot.Start()
}
