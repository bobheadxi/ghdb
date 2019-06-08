package log

import (
	"go.uber.org/zap"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
}

func NewDefaultLogger() Logger {
	return zap.NewExample().Sugar()
}
