package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

var LogLevels = []string{LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError}

/*
LocalLevelToTransform 本地日志等级
*/
func LocalLevelToTransform(localLevel string) zapcore.Level {
	switch localLevel {
	case LogLevelDebug:
		return zap.DebugLevel
	case LogLevelInfo:
		return zap.InfoLevel
	case LogLevelWarn:
		return zap.WarnLevel
	case LogLevelError:
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}

}
