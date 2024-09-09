package app

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/taskemapp/server/apps/notification/internal/config"
	"github.com/taskemapp/server/apps/notification/internal/notifier"
	notify "github.com/taskemapp/server/apps/notification/pkg/notifier"
	"github.com/taskemapp/server/libs/queue"
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

	fx.Provide(queue.NewConfig),
	fx.Provide(fx.Annotate(queue.NewMQ, fx.As(new(queue.Queue)))),
	fx.Provide(fx.Annotate(notifier.NewEmail, fx.As(new(notifier.EmailNotifier)))),

	fx.Invoke(
		func(logger *zap.Logger, c config.Config, mq queue.Queue, e notifier.EmailNotifier) {
			logger.Sugar().Info("Starting app: env - ", c.AppEnv)
			go func() {
				err := mq.Consume(
					notify.ChannelEmail,
					func(msg queue.Message) {
						var n notify.EmailNotification
						err := json.Unmarshal(msg.Body, &n)
						if err != nil {
							logger.Sugar().Errorf("Error unmarshalling email notification: %s", err)
						}

						logger.Sugar().Infof("Recieved message: %s", msg.Body)
						// TODO: temporary use a smtp user, because didn't have deployed smtp server
						// later need to change to a `n.From`
						email := notifier.EmailMsg{
							From:    c.SmtpUsername,
							To:      []string{n.To},
							Subject: n.Title,
							Body:    n.Message,
						}

						err = e.Send(email)
						if err != nil {
							logger.Sugar().Errorf("Error sending email: %s", err)
						}
					},
				)

				if err != nil {
					logger.Sugar().Errorf("Error consuming messages: %s", err)
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
	return amqp.Dial(c.RabbitMq.Url)
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
