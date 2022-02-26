package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type SimpleZapLoggerOption struct {
	LoggerFileMaxSize    int  `json:"logger_file_max_size" yaml:"logger_file_max_size"`
	LoggerFileMaxAge     int  `json:"logger_file_max_age" yaml:"logger_file_max_age"`
	LoggerFileMaxBackups int  `json:"logger_file_max_backups" yaml:"logger_file_max_backups"`
	LoggerFileCompress   bool `json:"logger_file_compress" yaml:"logger_file_compress"`

	LoggerFileLogDir      string `json:"logger_file_log_dir" yaml:"logger_file_log_dir"`
	LoggerFileLogName     string `json:"logger_file_log_name" yaml:"logger_file_log_name"`
	LoggerFileMinLogLevel string `json:"logger_file_min_log_level" yaml:"logger_file_min_log_level"`
	LoggerFileStdoutFlag  bool   `json:"logger_file_stdout_flag" yaml:"logger_file_stdout_flag"`
	LoggerFileCallDepth   int    `json:"logger_file_call_depth" yaml:"logger_file_call_depth"`

	LoggerFileMonoFile bool `json:"logger_file_mono_file" yaml:"logger_file_mono_file"`
}

type SimpleZapLogger struct {
	*zap.Logger

	SimpleZapLoggerOption
}

const (
	defaultLoggerFileMaxSize     = 256          // 默认 文件大小限制 256 MB
	defaultLoggerFileMaxAge      = 30           // 默认 日志文件保留天数 30天
	defaultLoggerFileMaxBackups  = 200          // 默认 最大保留日志文件数量 200
	defaultLoggerFileCompress    = false        // 默认 是否压缩处理 不压缩
	defaultLoggerFileLogDir      = "./log/"     // 默认 当前地址
	defaultLoggerFileLogName     = "logger.log" // 默认 当前地址
	defaultLoggerFileMinLogLevel = "debug"      // 默认 当前地址
	defaultLoggerFileStdoutFlag  = true         // 默认 同步输出到控制台
	defaultLoggerFileCallDepth   = 0            //默认 调用层级
	defaultLoggerFileMonoFile    = true         //默认 默认是单文件. 所有日志等级文件均在一个文件中
)

func NewSimpleZapLogger(options ...Option) *SimpleZapLogger {

	ret := &SimpleZapLogger{
		SimpleZapLoggerOption: SimpleZapLoggerOption{
			LoggerFileMaxSize:     defaultLoggerFileMaxSize,
			LoggerFileMaxAge:      defaultLoggerFileMaxAge,
			LoggerFileMaxBackups:  defaultLoggerFileMaxBackups,
			LoggerFileCompress:    defaultLoggerFileCompress,
			LoggerFileLogDir:      defaultLoggerFileLogDir,
			LoggerFileLogName:     defaultLoggerFileLogName,
			LoggerFileMinLogLevel: defaultLoggerFileMinLogLevel,
			LoggerFileStdoutFlag:  defaultLoggerFileStdoutFlag,
			LoggerFileCallDepth:   defaultLoggerFileCallDepth,
			LoggerFileMonoFile:    defaultLoggerFileMonoFile,
		},
	}

	for i := range options {
		options[i](ret)
	}
	if ret.LoggerFileMonoFile {
		ret.monoFile()
	} else {
		ret.multiFile()
	}
	return ret
}

func (s *SimpleZapLogger) Close() error {
	return s.Sync()
}

func (s *SimpleZapLogger) monoFile() {

	hook := lumberjack.Logger{
		Filename:   s.LoggerFileLogDir + s.LoggerFileLogName, // 日志文件路径
		MaxSize:    s.LoggerFileMaxSize,                      // 文件大小限制,单位MB, 默认 256MB
		MaxAge:     s.LoggerFileMaxAge,                       // 日志文件保留天数
		MaxBackups: s.LoggerFileMaxBackups,                   // 最大保留日志文件数量
		Compress:   s.LoggerFileCompress,                     // 是否压缩处理
		LocalTime:  true,
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "t",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "trace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	syncerList := []zapcore.WriteSyncer{
		zapcore.AddSync(&hook),
	}

	if s.LoggerFileStdoutFlag {
		// 同时输出到控制台
		syncerList = append(syncerList, zapcore.AddSync(os.Stdout))
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),          // 编码器配置
		zapcore.NewMultiWriteSyncer(syncerList...),     // 打印到控制台和文件
		LocalLevelToTransform(s.LoggerFileMinLogLevel), // 设置 最低输出级别
	)
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 增加调用者的层级
	skip := zap.AddCallerSkip(s.LoggerFileCallDepth)
	// 开启文件及行号
	development := zap.Development()
	// 构造日志
	s.Logger = zap.New(core, caller, skip, development)
}

func (s *SimpleZapLogger) multiFile() {
	var coreList []zapcore.Core
	minLevel := LocalLevelToTransform(s.LoggerFileMinLogLevel)
	for i := range LogLevels {
		curLevel := LocalLevelToTransform(LogLevels[i])
		if minLevel > curLevel {
			continue
		}
		coreList = append(coreList, s.zapCore(LogLevels[i]))
	}

	core := zapcore.NewTee(coreList...)
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 增加调用者的层级
	skip := zap.AddCallerSkip(s.LoggerFileCallDepth)
	// 开启文件及行号
	development := zap.Development()
	// 构造日志
	s.Logger = zap.New(core, caller, skip, development)
}

func (s *SimpleZapLogger) zapCore(level string) zapcore.Core {
	path := s.LoggerFileLogDir + s.LoggerFileLogName + "." + level
	hook := lumberjack.Logger{
		Filename:   path,                   // 日志文件路径
		MaxSize:    s.LoggerFileMaxSize,    // 文件大小限制,单位MB, 默认 256MB
		MaxAge:     s.LoggerFileMaxAge,     // 日志文件保留天数
		MaxBackups: s.LoggerFileMaxBackups, // 最大保留日志文件数量
		Compress:   s.LoggerFileCompress,   // 是否压缩处理
		LocalTime:  true,
	}

	syncerList := []zapcore.WriteSyncer{
		zapcore.AddSync(&hook),
	}

	if s.LoggerFileStdoutFlag {
		// 同时输出到控制台
		syncerList = append(syncerList, zapcore.AddSync(os.Stdout))
	}

	zapLevel := LocalLevelToTransform(level)

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "t",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "trace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}), // 编码器配置
		zapcore.NewMultiWriteSyncer(syncerList...), // 打印到控制台和文件
		zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl == zapLevel
		}), // 设置 最低输出级别
	)
}
