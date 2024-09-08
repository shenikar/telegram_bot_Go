package adapter

import (
	"log"
	"os"
	"strconv"
	"telegram_bot_go/domain"
	"telegram_bot_go/repository"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type TgBot struct {
	api         *tgbot.BotAPI
	hashService domain.HashWorder
	userRepo    *repository.UserRepo
}

func NewBot(api *tgbot.BotAPI, hashService domain.HashWorder, userRepo *repository.UserRepo) *TgBot {
	return &TgBot{
		api:         api,
		hashService: hashService,
		userRepo:    userRepo,
	}
}

func (b *TgBot) Start() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env: %v", err)
	}
	timeoutStr := os.Getenv("TIMEOUT_BOT")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		log.Fatalf("Invalid TIMEOUT_BOT value: %v", err)
	}

	u := tgbot.NewUpdate(0)
	u.Timeout = timeout
	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			userID := update.Message.From.ID
			text := update.Message.Text

			limitRequest, err := b.userRepo.LimitRequest(int(userID))
			if err != nil {
				log.Printf("Error getting limit request: %v", err)
			}

			if limitRequest {
				msg := tgbot.NewMessage(update.Message.Chat.ID, "Request limit. Please try again later.")
				b.api.Send(msg)
				continue
			}

			err = b.userRepo.SaveRequest(int(userID))
			if err != nil {
				log.Printf("Error saving request: %v", err)
				continue
			}

			if text == "/start" {
				msg := tgbot.NewMessage(update.Message.Chat.ID, "Hello! Please enter Md5 hash.")
				b.api.Send(msg)
				continue
			}
			// проверка, то что это хеш md5
			if len(text) == 32 {
				if originalWord, found := b.hashService.GetWordMulti(text); found {
					msg := tgbot.NewMessage(update.Message.Chat.ID, "Original word: "+originalWord)
					b.api.Send(msg)
				} else {
					msg := tgbot.NewMessage(update.Message.Chat.ID, "Hash not found")
					b.api.Send(msg)
				}

			} else {
				msg := tgbot.NewMessage(update.Message.Chat.ID, "Invalid hash")
				b.api.Send(msg)
			}
		}
	}

}
