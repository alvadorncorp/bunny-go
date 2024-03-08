package logger

type Logger interface {
	Flush() error
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(err error, msg string, fields ...Field)
	With(fields ...Field) Logger
}
