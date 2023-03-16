package logger

import (
	"go.uber.org/zap"
)

var Z *zap.Logger

func InitLogger() *zap.Logger {
	Z, _ = zap.NewProduction()
	return Z
}
