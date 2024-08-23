package broker

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

type Opts struct {
	fx.In
	Log  *zap.Logger
	Conn *amqp.Connection
}

type Mq struct {
	log  *zap.Logger
	conn *amqp.Connection
}

func New(opts Opts) *Mq {
	return &Mq{
		log:  opts.Log,
		conn: opts.Conn,
	}
}

func (m *Mq) Send(opts SendOpts, chName Channel) error {
	ch, err := m.conn.Channel()
	if err != nil {
		m.log.Sugar().Errorf("Failed open channel: %s", err)
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		string(chName),
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		m.log.Sugar().Errorf("Failed queue declare: %s", err)
		return err
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
			Headers:     opts.Headers,
			ContentType: opts.ContentType,
			Timestamp:   time.Now(),
			Body:        opts.Body,
		})
	if err != nil {
		m.log.Sugar().Errorf("Failed to send message: %s", err)
		return err
	}
	return nil
}
