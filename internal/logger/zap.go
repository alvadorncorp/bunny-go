package logger

import (
	"go.uber.org/zap"
)

type zapLogger struct {
	logger *zap.Logger
}

func NewZapLogger() (Logger, error) {
	l, err := zap.NewProduction()

	if err != nil {
		return nil, err
	}

	return &zapLogger{
		logger: l,
	}, nil
}

func (z *zapLogger) Flush() error {
	return z.logger.Sync()
}

func (z *zapLogger) Debug(msg string, fields ...Field) {
	z.logger.Debug(msg, fields...)
}

func (z *zapLogger) Info(msg string, fields ...Field) {
	z.logger.Info(msg, fields...)
}

func (z *zapLogger) Warn(msg string, fields ...Field) {
	z.logger.Warn(msg, fields...)
}

func (z *zapLogger) Error(err error, msg string, fields ...Field) {
	errField := zap.Error(err)
	allFields := append([]Field{errField}, fields...)

	z.logger.Error(msg, allFields...)
}

func (z *zapLogger) With(fields ...Field) Logger {
	newLogger := z.logger.With(fields...)
	return &zapLogger{
		logger: newLogger,
	}
}
