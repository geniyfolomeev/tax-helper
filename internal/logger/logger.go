package logger

import "go.uber.org/zap"

type Logger interface {
	Info(args ...any)
	Warn(args ...any)
	Error(args ...any)
	Errorf(template string, args ...interface{})
	Debug(args ...any)
	Fatal(args ...any)
}

func NewLogger(mode string) (*zap.SugaredLogger, error) {
	initLoggerFunc := zap.NewProduction
	if mode != "prod" {
		initLoggerFunc = zap.NewDevelopment
	}

	l, err := initLoggerFunc()
	if err != nil {
		return nil, err
	}
	return l.Sugar(), nil
}
