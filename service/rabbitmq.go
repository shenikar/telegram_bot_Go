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
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	log.Printf("Queue %s declared", q.Name) // Лог для проверки успешного создания очереди

	return &RabbitMQService{
		channel: ch,
		queue:   q,
	}, nil
}

func (r *RabbitMQService) CreateCallbackQueue() (amqp.Queue, error) {
	q, err := r.channel.QueueDeclare(
		"",
		false,
		true,
		true,
		false,
		nil,
	)
	if err != nil {
		return q, err
	}
	log.Printf("Callback queue %s created", q.Name) // Лог для проверки создания обратной очереди
	return q, nil
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
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	log.Printf("Consuming messages from %s", callbackQueueName)
	return msgs, nil
}
