package logger

type mockLogger struct {
}

func NewMockLogger() Logger {
	return &mockLogger{}
}

func (m mockLogger) Flush() error {
	return nil
}

func (m mockLogger) Debug(msg string, fields ...Field) {
	return
}

func (m mockLogger) Info(msg string, fields ...Field) {
	return
}

func (m mockLogger) Warn(msg string, fields ...Field) {
	return
}

func (m mockLogger) Error(err error, msg string, fields ...Field) {
	return
}

func (m mockLogger) With(fields ...Field) Logger {
	return m
}
