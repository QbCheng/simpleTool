package simpleZap

import "strings"

type Option func(*SimpleZapLogger)

// WithLoggerFileMaxSize 日志文件最大大小
func WithLoggerFileMaxSize(loggerFileMaxSize int) Option {
	return func(sz *SimpleZapLogger) {
		sz.LoggerFileMaxSize = loggerFileMaxSize
	}
}

// WithLoggerFileMaxAge 日志文件保留天数, 根据日志名中时间戳决定
func WithLoggerFileMaxAge(loggerFileMaxSize int) Option {
	return func(sz *SimpleZapLogger) {
		sz.LoggerFileMaxAge = loggerFileMaxSize
	}
}

// WithLoggerFileMaxBackups 最大保留日志文件数量
func WithLoggerFileMaxBackups(LoggerFileMaxBackups int) Option {
	return func(sz *SimpleZapLogger) {
		sz.LoggerFileMaxBackups = LoggerFileMaxBackups
	}
}

// WithLoggerFileCompress 是否压缩处理
func WithLoggerFileCompress(LoggerFileCompress bool) Option {
	return func(sz *SimpleZapLogger) {
		sz.LoggerFileCompress = LoggerFileCompress
	}
}

// WithLoggerFileLogDir 日志文件路径. 默认 ./test
func WithLoggerFileLogDir(LoggerFileLogDir string) Option {
	return func(sz *SimpleZapLogger) {
		if !strings.HasSuffix(LoggerFileLogDir, "/") {
			LoggerFileLogDir = LoggerFileLogDir + "/"
		}
		sz.LoggerFileLogDir = LoggerFileLogDir
	}
}

// WithLoggerFileLogName 日志文件名. 默认 ./logger.log
func WithLoggerFileLogName(LoggerFileLogName string) Option {
	return func(sz *SimpleZapLogger) {
		sz.LoggerFileLogName = LoggerFileLogName
	}
}

// WithLoggerFileMinLogLevel 最低输出级别
func WithLoggerFileMinLogLevel(LoggerFileMinLogLevel string) Option {
	return func(sz *SimpleZapLogger) {
		sz.LoggerFileMinLogLevel = LoggerFileMinLogLevel
	}
}

// WithLoggerFileStdoutFlag 同步输出到控制台
func WithLoggerFileStdoutFlag(LoggerFileStdoutFlag bool) Option {
	return func(sz *SimpleZapLogger) {
		sz.LoggerFileStdoutFlag = LoggerFileStdoutFlag
	}
}

// WithLoggerFileCallDepth 堆栈打印深度
func WithLoggerFileCallDepth(LoggerFileCallDepth int) Option {
	return func(sz *SimpleZapLogger) {
		sz.LoggerFileCallDepth = LoggerFileCallDepth
	}
}

// WithLoggerFileMonoFile 根据文件等级切分日志.
// 比如 WithLoggerFileLogPath 设置为 ./test, WithLoggerFileMonoFile 设置为 false. 会 输出 test.debug.log, test.info.log, test.warm.log, test.error.log
func WithLoggerFileMonoFile(loggerFileMonoFile bool) Option {
	return func(sz *SimpleZapLogger) {
		sz.LoggerFileMonoFile = loggerFileMonoFile
	}
}
