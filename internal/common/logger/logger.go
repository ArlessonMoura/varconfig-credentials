// Package logger provides centralized logging functionality
package logger

import (
	"context"
	"fmt"
	"log"
	"os"
)

// Level represents log levels
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Logger interface defines logging operations
type Logger interface {
	Debug(ctx context.Context, format string, args ...any)
	Info(ctx context.Context, format string, args ...any)
	Warn(ctx context.Context, format string, args ...any)
	Error(ctx context.Context, format string, args ...any)
}

// DefaultLogger implements a simple console logger
type DefaultLogger struct {
	level Level
}

// NewDefaultLogger creates a new default logger
func NewDefaultLogger(level Level) *DefaultLogger {
	return &DefaultLogger{level: level}
}

func (l *DefaultLogger) Debug(ctx context.Context, format string, args ...any) {
	if l.level <= LevelDebug {
		log.Printf("[DEBUG] %s", fmt.Sprintf(format, args...))
	}
}

func (l *DefaultLogger) Info(ctx context.Context, format string, args ...any) {
	if l.level <= LevelInfo {
		log.Printf("[INFO] %s", fmt.Sprintf(format, args...))
	}
}

func (l *DefaultLogger) Warn(ctx context.Context, format string, args ...any) {
	if l.level <= LevelWarn {
		log.Printf("[WARN] %s", fmt.Sprintf(format, args...))
	}
}

func (l *DefaultLogger) Error(ctx context.Context, format string, args ...any) {
	if l.level <= LevelError {
		log.Printf("[ERROR] %s", fmt.Sprintf(format, args...))
	}
}

// Global logger instance
var defaultLogger Logger = NewDefaultLogger(LevelInfo)

// SetGlobalLogger sets the global logger instance
func SetGlobalLogger(logger Logger) {
	defaultLogger = logger
}

// GetGlobalLogger returns the global logger instance
func GetGlobalLogger() Logger {
	return defaultLogger
}

// Debug logs a debug message using the global logger
func Debug(ctx context.Context, format string, args ...any) {
	defaultLogger.Debug(ctx, format, args...)
}

// Info logs an info message using the global logger
func Info(ctx context.Context, format string, args ...any) {
	defaultLogger.Info(ctx, format, args...)
}

// Warn logs a warning message using the global logger
func Warn(ctx context.Context, format string, args ...any) {
	defaultLogger.Warn(ctx, format, args...)
}

// Error logs an error message using the global logger
func Error(ctx context.Context, format string, args ...any) {
	defaultLogger.Error(ctx, format, args...)
}

// Initialize with environment variable
func init() {
	if os.Getenv("LOG_LEVEL") == "debug" {
		SetGlobalLogger(NewDefaultLogger(LevelDebug))
	}
}
