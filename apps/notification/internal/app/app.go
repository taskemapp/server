package app

import (
	"crypto/tls"
	"github.com/jordan-wright/email"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/taskemapp/server/apps/notification/internal/broker"
	"github.com/taskemapp/server/apps/notification/internal/config"
	"github.com/taskemapp/server/apps/notification/internal/notifier"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/smtp"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

var App = fx.Options(
	fx.Provide(setupConfig),
	fx.Provide(setupLogger),
	fx.Provide(setupRabbitMq),
	fx.Provide(setupSmtp),

	fx.Provide(broker.New),
	fx.Provide(fx.Annotate(notifier.NewEmail, fx.As(new(notifier.Notifier)))),

	fx.Invoke(
		func(logger *zap.Logger, c config.Config, mq *broker.Mq, e notifier.Notifier) {
			logger.Sugar().Info("Starting app: env - ", c.AppEnv)
			go func() {
				e.Send()
			}()
			go func() {

				err := mq.Receive("email")
				if err != nil {
					logger.Fatal("Error starting app", zap.Error(err))
				}
			}()
		},
	),
)

func setupConfig() (config.Config, error) {
	cfg, err := config.New()
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func setupRabbitMq(c config.Config) (*amqp.Connection, error) {
	return amqp.Dial(c.RabbitMqUrl)
}

func setupSmtp(c config.Config) (*email.Pool, error) {
	conf := &tls.Config{ServerName: c.SmtpHost}
	return email.NewPool(
		c.SmtpUrl,
		10,
		smtp.PlainAuth(
			"",
			c.SmtpUsername,
			c.SmtpPassword,
			c.SmtpHost,
		),
		conf,
	)
}

func setupLogger(c config.Config) *zap.Logger {
	var log *zap.Logger

	switch c.AppEnv {
	case envDev:
		log, _ = zap.NewDevelopment()
	case envProd:
		log, _ = zap.NewProduction()
	default:
		log, _ = zap.NewDevelopment()
	}

	return log
}
