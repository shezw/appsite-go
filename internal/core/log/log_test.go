// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package log

import (
	"context"
	"testing"
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
	SetLogger(mock)

	ctx := context.Background()

	Debug(ctx, "test debug")
	if mock.LastMsg != "test debug" || mock.LastLvl != "debug" {
		t.Errorf("Debug() failed")
	}

	Info(ctx, "test info")
	if mock.LastMsg != "test info" || mock.LastLvl != "info" {
		t.Errorf("Info() failed")
	}

	Warn(ctx, "test warn")
	if mock.LastMsg != "test warn" || mock.LastLvl != "warn" {
		t.Errorf("Warn() failed")
	}

	Error(ctx, "test error")
	if mock.LastMsg != "test error" || mock.LastLvl != "error" {
		t.Errorf("Error() failed")
	}
	
	Fatal(ctx, "test fatal")
	if mock.LastMsg != "test fatal" || mock.LastLvl != "fatal" {
		t.Errorf("Fatal() failed")
	}

	// Test GetLogger
	if GetLogger() != mock {
		t.Errorf("GetLogger() failed")
	}
}

func TestZapLogger(t *testing.T) {
	// Initialize ZapLogger
	logger, err := NewZapLogger("debug")
	if err != nil {
		t.Fatalf("NewZapLogger failed: %v", err)
	}

	// Just call methods to ensure no panic and coverage
	// Since it writes to stdout, verifying output programmatically is harder without redirecting Stdout or using zap test observer
	// For this phase, code coverage of having executed the lines is the goal.
	ctx := context.Background()
	logger.Debug(ctx, "zap debug", "key", "val")
	logger.Info(ctx, "zap info")
	logger.Warn(ctx, "zap warn")
	logger.Error(ctx, "zap error")
	// logger.Fatal call os.Exit, so we skip it or mock it
	
	// Test fields conversion
	fields := logger.toFields("key1", "val1", "key2") // Odd number
	if len(fields) != 1 {
		t.Errorf("toFields odd number failed")
	}
	
	if err := logger.Sync(); err != nil {
		// Sync might fail on stdout/stderr on some systems (inappropriate ioctl for device), 
		// but we just want to ensure method exists.
		// t.Logf("Sync returned error: %v", err)
	}
}
