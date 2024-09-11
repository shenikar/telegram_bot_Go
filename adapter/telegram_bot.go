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
	timeout := b.getTimeout()

	u := tgbot.NewUpdate(0)
	u.Timeout = timeout
	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			b.handleMessage(update.Message)
		}
	}
}

func (b *TgBot) getTimeout() int {
	timeoutStr := os.Getenv("TIMEOUT_BOT")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		log.Fatalf("Invalid TIMEOUT_BOT value: %v", err)
	}
	return timeout
}

func (b *TgBot) handleMessage(m *tgbot.Message) {
	userID := m.From.ID
	text := m.Text

	limitRequest, err := b.userRepo.LimitRequest(int(userID))
	if err != nil {
		log.Printf("Error getting limit request: %v", err)
		return
	}

	if limitRequest {
		msg := tgbot.NewMessage(m.Chat.ID, "Request limit. Please try again later.")
		b.api.Send(msg)
		return
	}

	err = b.userRepo.SaveRequest(int(userID))
	if err != nil {
		log.Printf("Error saving request: %v", err)
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
