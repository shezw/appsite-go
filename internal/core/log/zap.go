// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package log

import (
"context"
"os"

"go.uber.org/zap"
"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
logger *zap.Logger
}

// NewZapLogger creates a new zap logger
func NewZapLogger(level string, format string) (*ZapLogger, error) {
// Config
encoderConfig := zap.NewProductionEncoderConfig()
encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

// Level
var l zapcore.Level
if err := l.UnmarshalText([]byte(level)); err != nil {
l = zap.InfoLevel
}

var encoder zapcore.Encoder
if format == "console" {
encoder = zapcore.NewConsoleEncoder(encoderConfig)
} else {
encoder = zapcore.NewJSONEncoder(encoderConfig)
}

core := zapcore.NewCore(
encoder,
zapcore.AddSync(os.Stdout),
l,
)

logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)) // Skip 1 level to point to caller of wrapper

return &ZapLogger{
logger: logger,
}, nil
}

func (l *ZapLogger) toFields(fields ...interface{}) []zap.Field {
// Simple conversion: assume alternating string keys and values
// Ideally we want structured logging, but the interface takes ...interface{} for flexibility
// For better structured logging, we might change interface to take map or variadic structs
// But complying with standard ...interface{}:
if len(fields) == 0 {
return nil
}

zapFields := make([]zap.Field, 0, len(fields)/2)
for i := 0; i < len(fields); i += 2 {
if i+1 < len(fields) {
key, ok := fields[i].(string)
if !ok {
key = "unknown_key"
}
zapFields = append(zapFields, zap.Any(key, fields[i+1]))
}
}
return zapFields
}

func (l *ZapLogger) Debug(ctx context.Context, msg string, fields ...interface{}) {
l.logger.Debug(msg, l.toFields(fields...)...)
}

func (l *ZapLogger) Info(ctx context.Context, msg string, fields ...interface{}) {
l.logger.Info(msg, l.toFields(fields...)...)
}

func (l *ZapLogger) Warn(ctx context.Context, msg string, fields ...interface{}) {
l.logger.Warn(msg, l.toFields(fields...)...)
}

func (l *ZapLogger) Error(ctx context.Context, msg string, fields ...interface{}) {
l.logger.Error(msg, l.toFields(fields...)...)
}

func (l *ZapLogger) Fatal(ctx context.Context, msg string, fields ...interface{}) {
l.logger.Fatal(msg, l.toFields(fields...)...)
}

func (l *ZapLogger) Sync() error {
return l.logger.Sync()
}
