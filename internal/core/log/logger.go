// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package log

import (
	"context"
)

// Logger defines the interface for the application logger
type Logger interface {
	Debug(ctx context.Context, msg string, fields ...interface{})
	Info(ctx context.Context, msg string, fields ...interface{})
	Warn(ctx context.Context, msg string, fields ...interface{})
	Error(ctx context.Context, msg string, fields ...interface{})
	Fatal(ctx context.Context, msg string, fields ...interface{})
	
	// Sync flushes any buffered log entries
	Sync() error
}

var (
	globalLogger Logger
)

// SetLogger sets the global logger instance
func SetLogger(l Logger) {
	globalLogger = l
}

// GetLogger returns the global logger instance
// If not set, it returns a no-op logger (or should panic in strict mode, but here we return nil/noop)
func GetLogger() Logger {
	return globalLogger
}

// Helper functions using global logger

func Debug(ctx context.Context, msg string, fields ...interface{}) {
	if globalLogger != nil {
		globalLogger.Debug(ctx, msg, fields...)
	}
}

func Info(ctx context.Context, msg string, fields ...interface{}) {
	if globalLogger != nil {
		globalLogger.Info(ctx, msg, fields...)
	}
}

func Warn(ctx context.Context, msg string, fields ...interface{}) {
	if globalLogger != nil {
		globalLogger.Warn(ctx, msg, fields...)
	}
}

func Error(ctx context.Context, msg string, fields ...interface{}) {
	if globalLogger != nil {
		globalLogger.Error(ctx, msg, fields...)
	}
}

func Fatal(ctx context.Context, msg string, fields ...interface{}) {
	if globalLogger != nil {
		globalLogger.Fatal(ctx, msg, fields...)
	}
}
