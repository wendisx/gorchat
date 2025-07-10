package constant

import "go.uber.org/zap/zapcore"

type LogLevel zapcore.Level

// 日志格式
const (
	TIMESTAMP = "2006-01-02 15:04:05"
)

// 默认日志级别
const (
	DEBUG LogLevel = iota - 1
	INFO
	WARN
	ERROR
	DPANIC
	PANIC
	FATAL

	_maxLogLevel = FATAL
	INVALID      = _maxLogLevel + 1
)

// 日志字段
const (
	// status: "success"
	STATUS         = "status"
	STATUS_SUCCESS = "success"
	STATUS_FAIL    = "fail"
)
