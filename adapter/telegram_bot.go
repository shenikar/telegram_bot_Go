package adapter

import (
	"log"
	"telegram_bot_go/config"
	"telegram_bot_go/service"

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
		msg := tgbot.NewMessage(m.Chat.ID, "Attempt limit. Please try again later.")
		b.api.Send(msg)
		return
	}

	if text == "/stats" {

		callbackQueue, err := b.rabbitMQ.CreateCallbackQueue()
		if err != nil {
			msg := tgbot.NewMessage(m.Chat.ID, "Error creating callback queue.")
			b.api.Send(msg)
			return
		}

		err = b.rabbitMQ.SendQueueWithReply("/stats", callbackQueue.Name, int64(userID))
		if err != nil {
			msg := tgbot.NewMessage(m.Chat.ID, "Error sending stats request to queue.")
			b.api.Send(msg)
			return
		}

		go func() {
			msgs, err := b.rabbitMQ.ConsumeResults(callbackQueue.Name)
			if err != nil {
				log.Printf("Failed to register a consumer for stats: %v", err)
				return
			}
			b.ListenForResults(msgs, m.Chat.ID)
		}()
		return
	}

	// Создаем временную callback очередь
	callbackQueue, err := b.rabbitMQ.CreateCallbackQueue()
	if err != nil {
		msg := tgbot.NewMessage(m.Chat.ID, "Error creating callback queue")
		b.api.Send(msg)
		return
	}

	// Отправляем сообщение с указанием callback очереди
	if err := b.rabbitMQ.SendQueueWithReply(text, callbackQueue.Name, int64(userID)); err != nil {
		msg := tgbot.NewMessage(m.Chat.ID, "Error sending to queue")
		b.api.Send(msg)
		return
	}

	go func() {
		msgs, err := b.rabbitMQ.ConsumeResults(callbackQueue.Name)
		if err != nil {
			log.Printf("Failed to register a consumer: %v", err)
			return
		}
		b.ListenForResults(msgs, m.Chat.ID)
	}()
}

func (b *TgBot) ListenForResults(msgs <-chan amqp.Delivery, chatID int64) {
	for d := range msgs {
		result := string(d.Body)
		msg := tgbot.NewMessage(chatID, result)
		b.api.Send(msg)
	}
}

// func (b *TgBot) processCommand(m *tgbot.Message, text string) {
// 	if text == "/start" {
// 		msg := tgbot.NewMessage(m.Chat.ID, "Hello! Please enter Md5 hash.")
// 		b.api.Send(msg)
// 		return
// 	}
// 	// проверка, то что это хеш md5
// 	if len(text) == 32 {
// 		if originalWord, found := b.hashService.GetWord(text); found {
// 			msg := tgbot.NewMessage(m.Chat.ID, "Original word: "+originalWord)
// 			b.api.Send(msg)
// 		} else {
// 			msg := tgbot.NewMessage(m.Chat.ID, "Hash not found")
// 			b.api.Send(msg)
// 		}

// 	} else {
// 		msg := tgbot.NewMessage(m.Chat.ID, "Invalid hash")
// 		b.api.Send(msg)
// 	}
// }
