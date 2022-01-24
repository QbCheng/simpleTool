package customLogger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func NewZapLogger(logPath string, minLogLevel zapcore.Level, syncStdout bool) *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   logPath, // 日志文件路径
		MaxSize:    200,     // 文件大小限制,单位MB, 默认 100MB
		MaxBackups: 30,      // 最大保留日志文件数量
		Compress:   false,   // 是否压缩处理
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

	syncList := []zapcore.WriteSyncer{
		zapcore.AddSync(&hook),
	}

	if syncStdout {
		syncList = append(syncList, zapcore.AddSync(os.Stdout))
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),    // 编码器配置
		zapcore.NewMultiWriteSyncer(syncList...), // 打印到控制台和文件
		minLogLevel,                              // 设置 最低输出级别
	)
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()

	// 构造日志
	return zap.New(core, caller, development)
}
