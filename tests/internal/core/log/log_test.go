// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package log_test

import (
	"context"
	"testing"

	"appsite-go/internal/core/log"
)

type MockLogger struct {
	LastMsg string
	LastLvl string
}

func (m *MockLogger) Debug(ctx context.Context, msg string, fields ...interface{}) {
	m.LastMsg = msg
	m.LastLvl = "debug"
}
func (m *MockLogger) Info(ctx context.Context, msg string, fields ...interface{}) {
	m.LastMsg = msg
	m.LastLvl = "info"
}
func (m *MockLogger) Warn(ctx context.Context, msg string, fields ...interface{}) {
	m.LastMsg = msg
	m.LastLvl = "warn"
}
func (m *MockLogger) Error(ctx context.Context, msg string, fields ...interface{}) {
	m.LastMsg = msg
	m.LastLvl = "error"
}
func (m *MockLogger) Fatal(ctx context.Context, msg string, fields ...interface{}) {
	m.LastMsg = msg
	m.LastLvl = "fatal"
}
func (m *MockLogger) Sync() error { return nil }

func TestGlobalLogger(t *testing.T) {
	mock := &MockLogger{}
	log.SetLogger(mock)

	ctx := context.Background()

	log.Debug(ctx, "test debug")
	if mock.LastMsg != "test debug" || mock.LastLvl != "debug" {
		t.Errorf("Debug() failed")
	}

	log.Info(ctx, "test info")
	if mock.LastMsg != "test info" || mock.LastLvl != "info" {
		t.Errorf("Info() failed")
	}

	log.Warn(ctx, "test warn")
	if mock.LastMsg != "test warn" || mock.LastLvl != "warn" {
		t.Errorf("Warn() failed")
	}

	log.Error(ctx, "test error")
	if mock.LastMsg != "test error" || mock.LastLvl != "error" {
		t.Errorf("Error() failed")
	}
	
	log.Fatal(ctx, "test fatal")
	if mock.LastMsg != "test fatal" || mock.LastLvl != "fatal" {
		t.Errorf("Fatal() failed")
	}

	// Test GetLogger
	if log.GetLogger() != mock {
		t.Errorf("GetLogger() failed")
	}
}

func TestZapLogger(t *testing.T) {
	// Initialize ZapLogger
	logger, err := log.NewZapLogger("debug")
	if err != nil {
		t.Fatalf("NewZapLogger failed: %v", err)
	}

	// Just call methods to ensure no panic and coverage
	ctx := context.Background()
	logger.Debug(ctx, "zap debug", "key", "val")
	logger.Info(ctx, "zap info")
	logger.Warn(ctx, "zap warn")
	logger.Error(ctx, "zap error")
	
	if err := logger.Sync(); err != nil {
		// Ignore sync error in tests
	}
}
