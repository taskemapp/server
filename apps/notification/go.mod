module github.com/taskemapp/server/apps/notification

go 1.23.0

require (
	github.com/joho/godotenv v1.5.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/rabbitmq/amqp091-go v1.10.0
	github.com/taskemapp/server/libs/queue v0.0.0-20240906120508-8de35b503387
	github.com/wneessen/go-mail v0.4.4
	go.uber.org/fx v1.22.2
	go.uber.org/zap v1.27.0
)

require (
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	go.uber.org/dig v1.18.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
)
