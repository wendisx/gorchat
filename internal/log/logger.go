package log

import (
	"log"
	"time"

	"github.com/wendisx/gorchat/internal/constant"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 自定义日志类型
type LoggerConfig = zap.Config
type SimpleLogger = *zap.Logger
type Logger = *zap.SugaredLogger

func newTimeEncode(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func NewLoggerConfig() LoggerConfig {
	config := zap.NewProductionConfig()
	config.Development = true
	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.EncodeTime = newTimeEncode
	return config
}

// 创建日志器
func NewLogger(level constant.LogLevel) SimpleLogger {
	cfg := NewLoggerConfig()
	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("[init] -- (internal/logger) status: fail\n")
	} else {
		log.Printf("[init] -- (internal/logger) status: success\n")
	}
	return logger
}

// 内部核心日志处理
func Log(level constant.LogLevel, logger Logger, motion string, status int, ext map[string]any) {
	statusStr := constant.STATUS_FAIL
	if status == 0 {
		statusStr = constant.STATUS_SUCCESS
	}
	var kv []any
	kv = append(kv, constant.STATUS, statusStr)
	for k, v := range ext {
		kv = append(kv, k, v)
	}
	switch level {
	case constant.DEBUG:
		logger.Debugw(
			motion,
			kv...,
		)
	case constant.INFO:
		logger.Infow(
			motion,
			kv...,
		)
	case constant.WARN:
		logger.Warnw(
			motion,
			kv...,
		)
	case constant.ERROR:
		logger.Errorw(
			motion,
			kv...,
		)
	case constant.PANIC:
		logger.Panicw(
			motion,
			kv...,
		)
	case constant.FATAL:
		logger.Fatalw(
			motion,
			kv...,
		)
	default:
		log.Printf("[run] -- (internal/logger) status: fail\n")
	}
}

func Debug(logger Logger, motion string, ext map[string]any) {
	Log(constant.DEBUG, logger, motion, 0, ext)
}

func Info(logger Logger, motion string, ext map[string]any) {
	Log(constant.INFO, logger, motion, 0, ext)
}

func Warn(logger Logger, motion string, ext map[string]any) {
	Log(constant.WARN, logger, motion, 0, ext)
}

func Error(logger Logger, motion string, ext map[string]any) {
	Log(constant.ERROR, logger, motion, 1, ext)
}

func Dpanic(logger Logger, motion string, ext map[string]any) {
	Log(constant.DPANIC, logger, motion, 1, ext)
}

func Panic(logger Logger, motion string, ext map[string]any) {
	Log(constant.PANIC, logger, motion, 1, ext)
}

func Fatal(logger Logger, motion string, ext map[string]any) {
	Log(constant.FATAL, logger, motion, 1, ext)
}
