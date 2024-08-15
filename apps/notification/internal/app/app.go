package app

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"notification/internal/broker"
	"notification/internal/config"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

var App = fx.Options(
	fx.Provide(setupConfig),
	fx.Provide(setupLogger),
	fx.Provide(setupRabbitMq),

	fx.Provide(broker.New),

	fx.Invoke(
		func(logger *zap.Logger, c config.Config, mq *broker.Mq) {
			logger.Sugar().Info("Starting app: env - ", c.AppEnv)

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
