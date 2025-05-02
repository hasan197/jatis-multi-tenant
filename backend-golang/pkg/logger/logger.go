package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// Log adalah instance global dari logger
	Log *zap.Logger
)

// Config menyimpan konfigurasi untuk logger
type Config struct {
	LogLevel      string
	LogFilePath   string
	MaxSize       int
	MaxBackups    int
	MaxAge        int
	Compress      bool
	ConsoleOutput bool
}

// InitLogger menginisialisasi logger dengan konfigurasi yang diberikan
func InitLogger(cfg *Config) error {
	// Setup log rotation
	writer := &lumberjack.Logger{
		Filename:   cfg.LogFilePath,
		MaxSize:    cfg.MaxSize,    // megabytes
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,     // days
		Compress:   cfg.Compress,
	}

	// Buat encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Setup output
	var cores []zapcore.Core

	// File output
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(writer),
		getLogLevel(cfg.LogLevel),
	)
	cores = append(cores, fileCore)

	// Console output jika diaktifkan
	if cfg.ConsoleOutput {
		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			getLogLevel(cfg.LogLevel),
		)
		cores = append(cores, consoleCore)
	}

	// Buat core
	core := zapcore.NewTee(cores...)

	// Buat logger
	Log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return nil
}

// getLogLevel mengkonversi string level ke zapcore.Level
func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// WithContext menambahkan context ke logger
func WithContext(fields ...zap.Field) *zap.Logger {
	return Log.With(fields...)
}

// Sync memastikan semua log ditulis ke disk
func Sync() error {
	return Log.Sync()
} 