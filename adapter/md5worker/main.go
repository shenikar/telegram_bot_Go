package main

import (
	"log"

	"telegram_bot_go/config"
	"telegram_bot_go/service"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := amqp.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	queueName := cfg.RabbitMQ.Queue
	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	hashService := service.NewHashService()
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			hash := string(d.Body)
			originalWord, found := hashService.GetWordMulti(hash)

			response := ""
			if found {
				response = originalWord
			} else {
				response = "Original word not found for hash: " + hash
			}

			log.Printf("Sending response to %s: %s", d.ReplyTo, response)

			// Отправка результата обратно
			err := ch.Publish(
				"",
				d.ReplyTo,
				false,
				false,
				amqp.Publishing{
					ContentType:   "text/plain",
					Body:          []byte(response),
					CorrelationId: d.CorrelationId,
				},
			)
			if err != nil {
				log.Printf("Failed to publish a response: %v", err)
			}

			d.Ack(false)
		}
	}()

	log.Println(" [*] Waiting for hash processing requests.")
	<-forever
}
