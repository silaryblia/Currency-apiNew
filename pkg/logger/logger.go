package logger

import (
	"Currency-apiNew/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// Для обратной совместимости
func InitLogger(cfg *config.LoggerConfig) error {
	var zapConfig zap.Config

	if cfg.Development {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
	}

	// Устанавливаем уровень логирования
	switch cfg.Level {
	case "debug":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// Настраиваем кодировку
	if cfg.Encoding == "json" {
		zapConfig.Encoding = "json"
	} else {
		zapConfig.Encoding = "console"
	}

	zapConfig.EncoderConfig.TimeKey = "timestamp"
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	Logger, err = zapConfig.Build()
	if err != nil {
		return err
	}

	return nil
}

//func Sync() {
//	if Logger != nil {
//		Logger.Sync()
//	}
//}

//package logger
//
//import (
//	"go.uber.org/zap"
//	"go.uber.org/zap/zapcore"
//)
//
//var Logger *zap.Logger
//
//func InitLogger() error {
//	config := zap.NewDevelopmentConfig()
//	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
//	config.EncoderConfig.TimeKey = "timestamp"
//	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
//
//	var err error
//	Logger, err = config.Build()
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
