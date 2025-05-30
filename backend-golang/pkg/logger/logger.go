package logger

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/jatis/sample-stack-golang/internal/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// Log adalah instance global dari logger
	Log *logrus.Logger
)

// InitLogger menginisialisasi logger dengan konfigurasi yang diberikan
func InitLogger(cfg *config.Config) error {
	if cfg == nil {
		return nil
	}

	// Buat instance logger baru
	Log = logrus.New()

	// Setup log rotation
	writer := &lumberjack.Logger{
		Filename:   cfg.Logging.FilePath,
		MaxSize:    cfg.Logging.MaxSize,    // megabytes
		MaxBackups: cfg.Logging.MaxBackups,
		MaxAge:     cfg.Logging.MaxAge,     // days
		Compress:   cfg.Logging.Compress,
	}

	// Buat direktori log jika belum ada
	if err := os.MkdirAll(filepath.Dir(cfg.Logging.FilePath), 0755); err != nil {
		return err
	}

	// Setup formatter
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z",
	})

	// Setup level
	level, err := logrus.ParseLevel(cfg.Logging.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	Log.SetLevel(level)

	// Setup output
	if cfg.Logging.Output == "stdout" {
		// Multi writer untuk file dan console
		Log.SetOutput(writer)
		Log.AddHook(&ConsoleHook{
			Writer:    os.Stdout,
			LogLevels: logrus.AllLevels,
		})
	} else {
		// Hanya file output
		Log.SetOutput(writer)
	}

	return nil
}

// ConsoleHook adalah hook untuk menulis log ke console
type ConsoleHook struct {
	Writer    *os.File
	LogLevels []logrus.Level
}

// Levels mengembalikan level yang didukung oleh hook
func (hook *ConsoleHook) Levels() []logrus.Level {
	return hook.LogLevels
}

// Fire menulis log ke console
func (hook *ConsoleHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write([]byte(line))
	return err
}

// WithContext menambahkan context ke logger
func WithContext(fields map[string]interface{}) *logrus.Entry {
	return Log.WithFields(fields)
}

// Sync memastikan semua log ditulis ke disk
func Sync() error {
	return nil // Logrus tidak memerlukan sync
} 