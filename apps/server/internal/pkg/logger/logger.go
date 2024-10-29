package logger

import (
	"github.com/go-faster/errors"
	"go.uber.org/zap"
)

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)

	With(fields ...zap.Field) Logger
	Named(string) Logger
	WithScope(scope string) Logger
	WithMethod(method string) Logger
	WithOptions(opts ...zap.Option) Logger
	WithComponent(component string) Logger
}

type logger struct {
	*zap.Config
	*zap.Logger
}

func New(cfg *zap.Config, fields ...zap.Field) (Logger, error) {
	baseLogger, err := cfg.Build()
	if err != nil {
		return nil, errors.Wrap(err, "new logger")
	}

	baseLogger = baseLogger.With(fields...)

	defer func() { _ = baseLogger.Sync() }()

	return &logger{
		Config: cfg,
		Logger: baseLogger,
	}, nil
}

func (l *logger) Named(name string) Logger {
	return &logger{
		Config: l.Config,
		Logger: l.Logger.Named(name),
	}
}

func (l *logger) With(fields ...zap.Field) Logger {
	return &logger{
		Config: l.Config,
		Logger: l.Logger.With(fields...),
	}
}

func (l *logger) WithMethod(method string) Logger {
	return &logger{
		Config: l.Config,
		Logger: l.Logger.With(zap.String("method", method)),
	}
}

func (l *logger) WithScope(scope string) Logger {
	return &logger{
		Config: l.Config,
		Logger: l.Logger.With(zap.String("scope", scope)),
	}
}

func (l *logger) WithOptions(opts ...zap.Option) Logger {
	return &logger{
		Config: l.Config,
		Logger: l.Logger.WithOptions(opts...),
	}
}

func (l *logger) WithComponent(component string) Logger {
	return &logger{
		Config: l.Config,
		Logger: l.Logger.With(zap.String("component", component)),
	}
}
