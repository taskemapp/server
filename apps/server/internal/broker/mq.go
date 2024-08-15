package broker

import (
	"context"
	"github.com/go-faster/errors"
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

func (m *Mq) Send(opts SendOpts, chName string) error {
	ch, err := m.conn.Channel()
	if err != nil {
		m.log.Sugar().Error(errors.Wrap(err, "Failed open channel"))
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		chName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		m.log.Sugar().Error(errors.Wrap(err, "Failed queue declare"))
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
			UserId:      "asdasdasdasdasdasdsadsa",
			Body:        opts.Body,
		})
	if err != nil {
		m.log.Sugar().Error(errors.Wrap(err, "failed to send message"))
		return err
	}
	return nil
}
