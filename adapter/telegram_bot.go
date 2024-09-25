package adapter

import (
	"fmt"
	"log"
	"telegram_bot_go/config"
	"telegram_bot_go/service"
	"time"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	amqp "github.com/rabbitmq/amqp091-go"
)

type TgBot struct {
	api          *tgbot.BotAPI
	hashService  service.HashWorder
	userService  *service.UserService
	rabbitMQ     *service.RabbitMQService
	statsService *service.StatsService
	cfg          *config.Config
}

func NewBot(api *tgbot.BotAPI, hashService service.HashWorder, userService *service.UserService, rabbitMQ *service.RabbitMQService, statsService *service.StatsService, cfg *config.Config) *TgBot {
	return &TgBot{
		api:          api,
		hashService:  hashService,
		userService:  userService,
		rabbitMQ:     rabbitMQ,
		statsService: statsService,
		cfg:          cfg,
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
		msg := tgbot.NewMessage(m.Chat.ID, "Attempt limit reached. Please try again later.")
		b.api.Send(msg)
		return
	}

	if text == "/start" {
		msg := tgbot.NewMessage(m.Chat.ID, "Welcome to the MD5 Bot! Please send a hash for decoding.")
		b.api.Send(msg)
		return
	}

	if text == "/stats" {
		b.handleStatsCommand(m.Chat.ID, int(userID))
		return
	}

	if len(text) != 32 {
		msg := tgbot.NewMessage(m.Chat.ID, "Invalid hash format. Please send a valid MD5 hash (32 characters).")
		b.api.Send(msg)
		return
	}

	b.handleHashProcessing(text, m.Chat.ID, int64(userID))
}

func (b *TgBot) handleHashProcessing(hash string, chatID int64, userID int64) {
	callbackQueue, err := b.rabbitMQ.CreateCallbackQueue()
	if err != nil {
		msg := tgbot.NewMessage(chatID, "Error creating callback queue.")
		b.api.Send(msg)
		return
	}

	// Отправляем хеш в RabbitMQ для обработки
	err = b.rabbitMQ.SendQueueWithReply(hash, callbackQueue.Name, userID)
	if err != nil {
		msg := tgbot.NewMessage(chatID, "Error sending hash to queue.")
		b.api.Send(msg)
		return
	}

	// Ожидание ответа от воркера
	go func() {
		msgs, err := b.rabbitMQ.ConsumeResults(callbackQueue.Name)
		if err != nil {
			log.Printf("Failed to register a consumer: %v", err)
			return
		}
		b.ListenForResults(msgs, chatID)
	}()
}

func (b *TgBot) handleStatsCommand(chatID int64, userID int) {
	stats, err := b.statsService.GetStats(userID)
	if err != nil {
		msg := tgbot.NewMessage(chatID, "Error retrieving stats.")
		b.api.Send(msg)
		return
	}

	response := "Statistics:\n"
	for _, stat := range stats {
		response += fmt.Sprintf("Hash: %s - Result: %s at %s\n", stat.Hash, stat.Result, stat.AttemptTime.Format(time.RFC1123))
	}

	msg := tgbot.NewMessage(chatID, response)
	b.api.Send(msg)
}

func (b *TgBot) ListenForResults(msgs <-chan amqp.Delivery, chatID int64) {
	for d := range msgs {
		result := string(d.Body)
		msg := tgbot.NewMessage(chatID, result)
		b.api.Send(msg)
	}
}
