package adapter

import (
	"log"
	"telegram_bot_go/config"
	"telegram_bot_go/service"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgBot struct {
	api         *tgbot.BotAPI
	hashService service.HashWorder
	userService *service.UserService
	rabbitMQ    *service.RabbitMQService
	cfg         *config.Config
}

func NewBot(api *tgbot.BotAPI, hashService service.HashWorder, userService *service.UserService, rabbitMQ *service.RabbitMQService, cfg *config.Config) *TgBot {
	return &TgBot{
		api:         api,
		hashService: hashService,
		userService: userService,
		rabbitMQ:    rabbitMQ,
		cfg:         cfg,
	}
}

func (b *TgBot) Start() {
	u := tgbot.NewUpdate(0)
	u.Timeout = b.cfg.Telegram.Timeout
	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			b.handleMessage(update.Message)
		}
	}
}

func (b *TgBot) handleMessage(m *tgbot.Message) {
	userID := m.From.ID
	text := m.Text

	limitAttempt, err := b.userService.LimitAttempt(int(userID))
	if err != nil {
		log.Printf("Error getting limit attempt: %v", err)
		return
	}

	if limitAttempt {
		msg := tgbot.NewMessage(m.Chat.ID, "Attempt limit. Please try again later.")
		b.api.Send(msg)
		return
	}

	err = b.userService.SaveAttempt(int(userID), text)
	if err != nil {
		log.Printf("Error saving request: %v", err)
		return
	}

	if err := b.rabbitMQ.SendQueue(text); err != nil {
		msg := tgbot.NewMessage(m.Chat.ID, "Error sending to queue")
		b.api.Send(msg)
		return
	}

	b.processCommand(m, text)
}

func (b *TgBot) processCommand(m *tgbot.Message, text string) {
	if text == "/start" {
		msg := tgbot.NewMessage(m.Chat.ID, "Hello! Please enter Md5 hash.")
		b.api.Send(msg)
		return
	}
	// проверка, то что это хеш md5
	if len(text) == 32 {
		if originalWord, found := b.hashService.GetWord(text); found {
			msg := tgbot.NewMessage(m.Chat.ID, "Original word: "+originalWord)
			b.api.Send(msg)
		} else {
			msg := tgbot.NewMessage(m.Chat.ID, "Hash not found")
			b.api.Send(msg)
		}

	} else {
		msg := tgbot.NewMessage(m.Chat.ID, "Invalid hash")
		b.api.Send(msg)
	}
}
