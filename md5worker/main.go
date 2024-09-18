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
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}
	hashService := service.NewHashService()
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			hash := string(d.Body)
			log.Printf("Received a message: %s", hash)
			originalWord, found := hashService.GetWordMulti(hash)
			if found {
				log.Printf("Found original word: %s", originalWord)
			} else {
				log.Printf("Original word not found for hash: %s", hash)
			}

		}
	}()
	log.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
