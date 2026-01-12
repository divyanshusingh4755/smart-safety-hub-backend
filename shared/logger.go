package shared

import "go.uber.org/zap"

type Logger = *zap.Logger

func NewLogger() Logger {
	l, _ := zap.NewProduction()
	return l
}
