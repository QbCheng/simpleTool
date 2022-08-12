package logger

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestNewSimpleZapLogger(t *testing.T) {
	logger := NewSimpleZapLogger(WithLoggerFileCallDepth(0))
	for i := 0; i < 10000; i++ {
		logger.Error("aa", zap.Any("test", map[string]interface{}{
			"a": "a",
			"b": "b",
			"c": "c",
			"d": "d",
		}))
		logger.Info("aa", zap.Any("test", map[string]interface{}{
			"a": "a",
			"b": "b",
			"c": "c",
			"d": "d",
		}))
		logger.Debug("aa", zap.Any("test", map[string]interface{}{
			"a": "a",
			"b": "b",
			"c": "c",
			"d": "d",
		}))

		logger.Warn("aa", zap.Any("test", map[string]interface{}{
			"a": "a",
			"b": "b",
			"c": "c",
			"d": "d",
		}))
	}
	err := logger.Close()
	assert.NoError(t, err)
}

func TestNewSimpleZapLoggerToLevel(t *testing.T) {
	logger := NewSimpleZapLogger(
		WithLoggerFileCallDepth(0),
		WithLoggerFileStdoutFlag(false),
		WithLoggerFileMaxSize(1))
	for i := 0; i < 100000; i++ {
		logger.Error("aa", zap.Any("test", map[string]interface{}{
			"a": "a",
			"b": "b",
			"c": "c",
			"d": "d",
		}))
		logger.Info("aa", zap.Any("test", map[string]interface{}{
			"a": "a",
			"b": "b",
			"c": "c",
			"d": "d",
		}))
		logger.Debug("aa", zap.Any("test", map[string]interface{}{
			"a": "a",
			"b": "b",
			"c": "c",
			"d": "d",
		}))

		logger.Warn("aa", zap.Any("test", map[string]interface{}{
			"a": "a",
			"b": "b",
			"c": "c",
			"d": "d",
		}))
	}
	err := logger.Close()
	assert.NoError(t, err)
}

func TestNewSimpleZapLoggerToMonoFile(t *testing.T) {
	logger := NewSimpleZapLogger(
		WithLoggerFileCallDepth(0),
		WithLoggerFileStdoutFlag(false),
		WithLoggerFileMaxSize(1),
		WithLoggerFileMonoFile(false),
	)
	for i := 0; i < 100000; i++ {
		logger.Error("aa", zap.Any("test", map[string]interface{}{
			"a": "a",
			"b": "b",
			"c": "c",
			"d": "d",
		}))
		logger.Info("aa", zap.Any("test", map[string]interface{}{
			"a": "a",
			"b": "b",
			"c": "c",
			"d": "d",
		}))
		logger.Debug("aa", zap.Any("test", map[string]interface{}{
			"a": "a",
			"b": "b",
			"c": "c",
			"d": "d",
		}))

		logger.Warn("aa", zap.Any("test", map[string]interface{}{
			"a": "a",
			"b": "b",
			"c": "c",
			"d": "d",
		}))
	}
	err := logger.Close()
	assert.NoError(t, err)
}

func TestNewSimpleZapLoggerToMonoFileAndError(t *testing.T) {
	logger := NewSimpleZapLogger(
		WithLoggerFileCallDepth(0),
		WithLoggerFileStdoutFlag(false),
		WithLoggerFileMaxSize(1),
		WithLoggerFileMonoFile(false),
		WithLoggerFileMinLogLevel(LogLevelError),
	)
	for i := 0; i < 100000; i++ {
		logger.Error("aa", zap.Any("test", map[string]interface{}{
			"a": "a",
			"b": "b",
			"c": "c",
			"d": "d",
		}))
		logger.Info("aa", zap.Any("test", map[string]interface{}{
			"a": "a",
			"b": "b",
			"c": "c",
			"d": "d",
		}))
		logger.Debug("aa", zap.Any("test", map[string]interface{}{
			"a": "a",
			"b": "b",
			"c": "c",
			"d": "d",
		}))

		logger.Warn("aa", zap.Any("test", map[string]interface{}{
			"a": "a",
			"b": "b",
			"c": "c",
			"d": "d",
		}))
	}
	err := logger.Close()
	assert.NoError(t, err)
}

func TestNewSimpleZapLoggerToDirAndName(t *testing.T) {
	logger := NewSimpleZapLogger(
		WithLoggerFileCallDepth(0),
		WithLoggerFileStdoutFlag(false),
		WithLoggerFileMaxSize(1),
		WithLoggerFileMonoFile(false),
		WithLoggerFileLogDir("./customLogger"),
		WithLoggerFileLogName("./customLoggerLogger.log"),
	)
	for i := 0; i < 100000; i++ {
		logger.Error("aa", zap.Any("test", map[string]interface{}{
			"a": "a",
			"b": "b",
			"c": "c",
			"d": "d",
		}))
		logger.Info("aa", zap.Any("test", map[string]interface{}{
			"a": "a",
			"b": "b",
			"c": "c",
			"d": "d",
		}))
		logger.Debug("aa", zap.Any("test", map[string]interface{}{
			"a": "a",
			"b": "b",
			"c": "c",
			"d": "d",
		}))

		logger.Warn("aa", zap.Any("test", map[string]interface{}{
			"a": "a",
			"b": "b",
			"c": "c",
			"d": "d",
		}))
	}
	err := logger.Close()
	assert.NoError(t, err)
}
