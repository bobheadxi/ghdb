package log

import (
	"go.uber.org/zap"
)

// Logger defines the interface used internally for logging
type Logger interface {
	Debug(args ...interface{})
	Debugf(msg string, args ...interface{})

	Info(args ...interface{})
	Infof(msg string, args ...interface{})

	Error(args ...interface{})
	Errorf(msg string, args ...interface{})
}

// NewDefaultLogger provides a basic Logger implementation backed by uber-go/zap
func NewDefaultLogger() Logger {
	return zap.NewExample().Sugar()
}
