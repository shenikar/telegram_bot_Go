package service

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQService struct {
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewRabbitMqService(amqpURI, queueName string) (*RabbitMQService, error) {
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}
	return &RabbitMQService{
		channel: ch,
		queue:   q,
	}, nil
}

func (r *RabbitMQService) SendQueue(hash string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := r.channel.PublishWithContext(
		ctx,          // Передача контекста
		"",           // Exchange
		r.queue.Name, // Queue name
		false,        // Mandatory
		false,        // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(hash),
		},
	)
	if err != nil {
		log.Printf("Failed to publish a message: %v", err)
		return err
	}
	log.Printf(" [x] Sent %s\n", hash)
	return nil
}
