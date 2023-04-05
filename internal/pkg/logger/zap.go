package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func createLogger(debugMode bool, logEncoder string) *zap.SugaredLogger {
	var options []zap.Option
	var cfg zapcore.EncoderConfig
	if debugMode {
		cfg = zap.NewDevelopmentEncoderConfig()
		options = append(options, zap.AddCaller())
	} else {
		cfg = zap.NewProductionEncoderConfig()
	}

	var encoder zapcore.Encoder
	if logEncoder == "json" {
		encoder = zapcore.NewJSONEncoder(cfg)
	} else {
		encoder = zapcore.NewConsoleEncoder(cfg)
	}

	var logLevel zapcore.Level
	if debugMode {
		logLevel = zapcore.DebugLevel
	} else {
		logLevel = zapcore.InfoLevel
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.Lock(zapcore.AddSync(os.Stdout)),
		logLevel,
	)

	return zap.New(core, options...).Sugar()
}

func CreateLocalLogger() *zap.SugaredLogger {
	return createLogger(true, "plain")
}

func CreateDevLogger() *zap.SugaredLogger {
	return createLogger(true, "json")
}

func CreateProdLogger() *zap.SugaredLogger {
	return createLogger(false, "json")
}
