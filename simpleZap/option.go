package logger

type Option func(*SimpleZapLogger)

func WithLoggerFileMaxSize(loggerFileMaxSize int) Option {
	return func(sz *SimpleZapLogger) {
		sz.LoggerFileMaxSize = loggerFileMaxSize
	}
}

func WithLoggerFileMaxAge(loggerFileMaxSize int) Option {
	return func(sz *SimpleZapLogger) {
		sz.LoggerFileMaxAge = loggerFileMaxSize
	}
}

func WithLoggerFileMaxBackups(LoggerFileMaxBackups int) Option {
	return func(sz *SimpleZapLogger) {
		sz.LoggerFileMaxBackups = LoggerFileMaxBackups
	}
}

func WithLoggerFileCompress(LoggerFileCompress bool) Option {
	return func(sz *SimpleZapLogger) {
		sz.LoggerFileCompress = LoggerFileCompress
	}
}

func WithLoggerFileLogPath(LoggerFileLogPath string) Option {
	return func(sz *SimpleZapLogger) {
		sz.LoggerFileLogPath = LoggerFileLogPath
	}
}

func WithLoggerFileMinLogLevel(LoggerFileMinLogLevel string) Option {
	return func(sz *SimpleZapLogger) {
		sz.LoggerFileMinLogLevel = LoggerFileMinLogLevel
	}
}

func WithLoggerFileStdoutFlag(LoggerFileStdoutFlag bool) Option {
	return func(sz *SimpleZapLogger) {
		sz.LoggerFileStdoutFlag = LoggerFileStdoutFlag
	}
}

func WithLoggerFileCallDepth(LoggerFileCallDepth int) Option {
	return func(sz *SimpleZapLogger) {
		sz.LoggerFileCallDepth = LoggerFileCallDepth
	}
}

func WithLoggerFileMonoFile(loggerFileMonoFile bool) Option {
	return func(sz *SimpleZapLogger) {
		sz.LoggerFileMonoFile = loggerFileMonoFile
	}
}
