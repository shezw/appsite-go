// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package operation_test

import (
	"testing"

	"appsite-go/internal/services/access/operation"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestAudit_Flow(t *testing.T) {
	db := setupDB(t)
	svc := operation.NewService(db)

	logEntry := &operation.AuditLog{
		UserID: "u_1",
		Action: "login",
		IP:     "127.0.0.1",
		Status: 200,
	}

	// 1. Record
	if err := svc.Record(logEntry); err != nil {
		t.Fatalf("Record failed: %v", err)
	}
	if logEntry.ID == "" {
		t.Error("ID not generated")
	}

	// 2. Find
	logs, err := svc.FindByUser("u_1", 10)
	if err != nil {
		t.Fatalf("Find failed: %v", err)
	}
	if len(logs) != 1 {
		t.Errorf("Expected 1 log, got %d", len(logs))
	}
	if logs[0].Action != "login" {
		t.Errorf("Mismatch action")
	}
}
