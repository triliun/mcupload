package logger

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	globalLogger *zap.Logger
	once         sync.Once
)

func Init(config Config) error {
	var err error
	once.Do(func() {
		globalLogger, err = createLogger(config)
	})

	return err
}

// createLogger creates a new zap logger based on configuration
func createLogger(config Config) (*zap.Logger, error) {
	var level zapcore.Level
	err := level.UnmarshalText([]byte(config.Level))
	if err != nil {
		level = zapcore.InfoLevel
	}

	encoderConfig := getEncoderConfig(config.Environment)

	// Create cores
	cores := []zapcore.Core{}

	// Console output for development
	if config.Environment == "development" {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.Lock(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	// File outpur for production
	if config.Environment == "production" {
		// Main log file
		var mainWriter *os.File
		mainWriter, err = getLogWriter(config.OutputPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create main log writer: %w", err)
		}

		var errorWriter *os.File
		errorWriter, err = getLogWriter(config.ErrorPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create error log writer: %w", err)
		}

		jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

		// Main log file
		mainCore := zapcore.NewCore(
			jsonEncoder,
			zapcore.AddSync(mainWriter),
			level,
		)

		// Error log file (only errors)
		errorCore := zapcore.NewCore(
			jsonEncoder,
			zapcore.AddSync(errorWriter),
			level,
		)

		cores = append(cores, mainCore, errorCore)
	}

	// Create the core
	core := zapcore.NewTee(cores...)

	// Add caller information in development

	var options []zap.Option
	if config.Environment == "development" {
		options = append(options, zap.AddCaller())
	}

	options = append(options,
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.Fields(
			zap.String("service", "api"),
			zap.String("environment", config.Environment),
		),
	)

	return zap.New(core, options...), nil
}

func getEncoderConfig(env string) zapcore.EncoderConfig {
	if env == "production" {
		return zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
	}

	// Development config
	return zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func getLogWriter(path string) (*os.File, error) {
	if path == "stdout" {
		return os.Stdout, nil
	}
	if path == "stderr" {
		return os.Stderr, nil
	}

	// Create directory if it doesn't exist
	dir := path[:strings.LastIndex(path, "/")]
	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		return nil, err
	}

	return os.OpenFile(dir, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
}

// Get returns the globalLogger instance
func Get() *zap.Logger {
	if globalLogger == nil {
		// Fallback to development if not initialized
		globalLogger, _ = createLogger(DefaultConfig())
	}
	return globalLogger
}

// Sync flushes any buffered log entries
func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

// Helper methods for common log operations
func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Get().Fatal(msg, fields...)
}

// FieldsWhithRequestID creates a new logger with request_id and additional fields
func FieldsWhithRequestID(ctx context.Context, fields ...zap.Field) *zap.Logger {
	requestID := zap.String("request_id", middleware.GetReqID(ctx))
	fields = append([]zap.Field{requestID}, fields...)

	return WhithFields(fields...)
}

// WhithFields creates a new logger with additional fields
func WhithFields(fields ...zap.Field) *zap.Logger {
	return Get().With(fields...)
}
