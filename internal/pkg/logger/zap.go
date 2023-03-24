package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Z *zap.SugaredLogger

// by default, it will init simple logger
func init() {
	InitLogger(true, "plain")
}

func InitLogger(debugMode bool, logEncoder string) {
	cfg := zap.NewProductionEncoderConfig()
	cfg.MessageKey = "msg"
	cfg.LevelKey = "level"
	cfg.TimeKey = "time"
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	core := zapcore.NewCore(
		newEncoder(logEncoder, cfg),
		zapcore.Lock(zapcore.AddSync(os.Stdout)),
		logLevel(debugMode),
	)

	logger := zap.New(core)

	Z = logger.Sugar()
}

func newEncoder(logEncoder string, cfg zapcore.EncoderConfig) zapcore.Encoder {
	if logEncoder == "json" {
		return zapcore.NewJSONEncoder(cfg)
	}
	return zapcore.NewConsoleEncoder(cfg)
}

func logLevel(debugMode bool) zapcore.Level {
	if debugMode {
		return zapcore.DebugLevel
	}
	return zapcore.InfoLevel
}
