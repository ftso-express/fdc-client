package logger

import (
	"errors"
	"log"
	"os"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	timeFormat = "[01-02|15:04:05.000]"
)

var (
	sugaredLogger *zap.SugaredLogger
)

func init() {
	sugaredLogger = createSugaredLogger(DefaultLoggerConfig())
}

type LoggerConfig struct {
	Level       string `toml:"level"` // valid values are: DEBUG, INFO, WARN, ERROR, DPANIC, PANIC, FATAL (zap)
	File        string `toml:"file"`
	MaxFileSize int    `toml:"max_file_size"` // In megabytes
	Console     bool   `toml:"console"`
}

func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Level:   "DEBUG",
		Console: true,
	}
}

func SetLogger(config LoggerConfig) {
	createSugaredLogger(config)
}

func createSugaredLogger(config LoggerConfig) *zap.SugaredLogger {
	atom := zap.NewAtomicLevel()
	cores := make([]zapcore.Core, 0)
	if config.Console {
		cores = append(cores, createConsoleLoggerCore(atom))
	}
	if len(config.File) > 0 {
		cores = append(cores, createFileLoggerCore(config, atom))
	}

	core := zapcore.NewTee(cores...)
	logger := zap.New(core,
		zap.AddStacktrace(zap.ErrorLevel),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
	defer func() {
		err := logger.Sync()
		if err != nil && !errors.Is(err, syscall.ENOTTY) && !errors.Is(err, syscall.EBADF) {
			log.Print("Failed to sync logger", err)
		}
	}()

	sugaredLogger = logger.Sugar()

	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		sugaredLogger.Errorf("Wrong level %s", config.Level)
	}
	atom.SetLevel(level)
	return sugaredLogger
}

func createFileLoggerCore(config LoggerConfig, atom zap.AtomicLevel) zapcore.Core {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename: config.File,
		MaxSize:  config.MaxFileSize,
	})
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeLevel = fileLevelEncoder
	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout(timeFormat)
	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		w,
		atom,
	)
}

func createConsoleLoggerCore(atom zap.AtomicLevel) zapcore.Core {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeLevel = consoleColorLevelEncoder
	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout(timeFormat)
	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		atom,
	)
}

func consoleColorLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	s, ok := levelToCapitalColorString[l]
	if !ok {
		s = unknownLevelColor.Wrap(l.CapitalString())
	}
	enc.AppendString(s)
}

func fileLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(l.CapitalString())
}

func Warnf(msg string, args ...interface{}) {
	sugaredLogger.Warnf(msg, args...)
}

func Errorf(msg string, args ...interface{}) {
	sugaredLogger.Errorf(msg, args...)
}

func Infof(msg string, args ...interface{}) {
	sugaredLogger.Infof(msg, args...)
}

func Debugf(msg string, args ...interface{}) {
	sugaredLogger.Debugf(msg, args...)
}

func Fatalf(msg string, args ...interface{}) {
	sugaredLogger.Fatalf(msg, args...)
}

func Panicf(msg string, args ...interface{}) {
	sugaredLogger.Panicf(msg, args...)
}

func Warn(args ...interface{}) {
	sugaredLogger.Warn(args...)
}

func Error(args ...interface{}) {
	sugaredLogger.Error(args...)
}

func Info(args ...interface{}) {
	sugaredLogger.Info(args...)
}

func Debug(args ...interface{}) {
	sugaredLogger.Debug(args...)
}

func Fatal(args ...interface{}) {
	sugaredLogger.Fatal(args...)
}

func Panic(args ...interface{}) {
	sugaredLogger.Panic(args...)
}
