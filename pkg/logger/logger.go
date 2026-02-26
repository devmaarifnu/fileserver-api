package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log *logrus.Logger

// Config holds logger configuration
type Config struct {
	Level      string
	Format     string
	Output     string
	FilePath   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

// Initialize sets up the logger
func Initialize(config Config) error {
	log = logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)

	// Set log format
	if config.Format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// Set output
	var writers []io.Writer

	switch config.Output {
	case "console":
		writers = append(writers, os.Stdout)
	case "file":
		fileWriter, err := getFileWriter(config)
		if err != nil {
			return err
		}
		writers = append(writers, fileWriter)
	case "both":
		fileWriter, err := getFileWriter(config)
		if err != nil {
			return err
		}
		writers = append(writers, os.Stdout, fileWriter)
	default:
		writers = append(writers, os.Stdout)
	}

	if len(writers) > 1 {
		log.SetOutput(io.MultiWriter(writers...))
	} else {
		log.SetOutput(writers[0])
	}

	return nil
}

// getFileWriter creates a file writer with rotation
func getFileWriter(config Config) (io.Writer, error) {
	// Ensure log directory exists
	logDir := filepath.Dir(config.FilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	return &lumberjack.Logger{
		Filename:   config.FilePath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}, nil
}

// GetLogger returns the logger instance
func GetLogger() *logrus.Logger {
	if log == nil {
		log = logrus.New()
	}
	return log
}

// WithFields creates a new entry with fields
func WithFields(fields logrus.Fields) *logrus.Entry {
	return GetLogger().WithFields(fields)
}

// WithField creates a new entry with a single field
func WithField(key string, value interface{}) *logrus.Entry {
	return GetLogger().WithField(key, value)
}

// Debug logs a debug message
func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

// Debugf logs a formatted debug message
func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

// Info logs an info message
func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

// Infof logs a formatted info message
func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

// Warn logs a warning message
func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

// Warnf logs a formatted warning message
func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

// Error logs an error message
func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

// Errorf logs a formatted error message
func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

// Fatal logs a fatal message and exits
func Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}

// Fatalf logs a formatted fatal message and exits
func Fatalf(format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}
