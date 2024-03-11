package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(lvl zapcore.Level) (Logger, error) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		timeStr := t.UTC().Format(time.RFC3339)
		encoder.AppendString(timeStr)
	}
	jsonCfg := zapcore.NewJSONEncoder(encoderCfg)
	ws := zapcore.Lock(os.Stdout)
	atom := zap.NewAtomicLevelAt(lvl)

	l := zap.New(zapcore.NewCore(jsonCfg, ws, atom))
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
