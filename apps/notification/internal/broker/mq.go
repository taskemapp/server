package broker

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"log"
	"sync"
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
	wg   sync.WaitGroup
}

func New(opts Opts) *Mq {
	return &Mq{
		log:  opts.Log,
		conn: opts.Conn,
		wg:   sync.WaitGroup{},
	}
}

func (m *Mq) Receive(queueName string) error {
	ch, err := m.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	println(fmt.Sprintf("%v:%v", queueName, time.Now()))

	consumer, err := ch.Consume(
		q.Name,
		fmt.Sprintf("%v:%v", queueName, time.Now()),
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	m.wg.Add(1)
	go func() {
		m.log.Sugar().Info("Start listen")
		defer m.wg.Done()
		for msg := range consumer {
			log.Printf("Received a message: %s", msg.Body)
		}
		m.log.Sugar().Info("Done listen")

	}()

	m.wg.Wait()

	return err
}
