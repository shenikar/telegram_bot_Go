package adapter

import (
	"log"
	"os"
	"strconv"
	"telegram_bot_go/domain"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type TgBot struct {
	api         *tgbot.BotAPI
	hashService domain.HashWorder
}

func NewBot(api *tgbot.BotAPI, hashService domain.HashWorder) *TgBot {
	return &TgBot{
		api:         api,
		hashService: hashService,
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
			text := update.Message.Text
			// проверка, то что это хеш md5
			if len(text) == 32 {
				if originalWord, found := b.hashService.GetWord(text); found {
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
