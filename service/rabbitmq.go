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

func (r *RabbitMQService) CreateCallbackQueue() (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		"",
		false,
		true,
		true,
		false,
		nil,
	)
}

func (r *RabbitMQService) SendQueueWithReply(hash string, callbackQueueName string, chatID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.channel.PublishWithContext(
		ctx,
		"",
		r.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(hash),
			ReplyTo:     callbackQueueName,
		},
	)
	if err != nil {
		log.Printf("Failed to publish a message: %v", err)
		return err
	}
	log.Printf(" [x] Sent %s to queue %s\n", hash, r.queue.Name)
	return nil
}

func (r *RabbitMQService) ConsumeResults(callbackQueueName string) (<-chan amqp.Delivery, error) {
	msgs, err := r.channel.Consume(
		callbackQueueName,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}
