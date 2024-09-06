package queue

import (
	"context"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"time"
)

type RabbitMQ struct {
	conn *amqp.Connection
	c    Config
}

func NewMQ(conn *amqp.Connection, c Config) *RabbitMQ {
	return &RabbitMQ{conn: conn, c: c}
}

func (r *RabbitMQ) Consume(name string, handler ConsumeFn) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to declare a queue")
	}

	consumer, err := ch.Consume(
		q.Name,
		"notification service",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return errors.Wrap(err, "Failed to consume")
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range consumer {
			handler(Message{
				ContentType: msg.ContentType,
				Body:        msg.Body,
			})
		}
	}()

	wg.Wait()

	return nil
}

func (r *RabbitMQ) Close() error {
	return r.conn.Close()
}

func (r *RabbitMQ) Publish(queue string, message Message) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to declare a queue")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: message.ContentType,
			Timestamp:   time.Now(),
			Body:        message.Body,
		})
	if err != nil {
		return errors.Wrap(err, "Failed to publish a message")
	}

	return nil
}
