package app

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/taskemapp/server/apps/notification/internal/broker"
	"github.com/taskemapp/server/apps/notification/internal/config"
	"github.com/taskemapp/server/apps/notification/internal/notifier"
	"github.com/wneessen/go-mail"
	"go.uber.org/fx"
	"go.uber.org/zap"
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
	fx.Provide(fx.Annotate(notifier.NewEmail, fx.As(new(notifier.EmailNotifier)))),

	fx.Invoke(
		func(logger *zap.Logger, c config.Config, mq *broker.Mq, e notifier.EmailNotifier) {
			logger.Sugar().Info("Starting app: env - ", c.AppEnv)
			go func() {
				err := mq.Receive("notifications")
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

func setupSmtp(c config.Config) (*mail.Client, error) {
	return mail.NewClient(c.SmtpHost,
		mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(c.SmtpUsername), mail.WithPassword(c.SmtpPassword),
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
