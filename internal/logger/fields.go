package logger

import (
	"go.uber.org/zap"
)

type Field = zap.Field

func String(key, value string) Field {
	return zap.String(key, value)
}
