// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package permission_test

import (
	"os"
	"testing"

	"appsite-go/internal/services/access/permission"
)

func TestCasbin_Enforce(t *testing.T) {
	// Standard Setup
	svc, err := permission.NewService("model.conf", "policy.csv")
	if err != nil {
		t.Fatalf("Failed to init service: %v", err)
	}

	tests := []struct {
		Sub      string
		Obj      string
		Act      string
		Expected bool
	}{
		{"alice", "/api/admin/dashboard", "read", true},   // alice is admin
		{"alice", "/api/admin/users", "write", true},      // alice is admin
		{"bob", "/api/admin/dashboard", "read", false},    // bob is member
		{"bob", "/api/user/profile", "read", true},        // bob is member
		{"charlie", "/api/user/profile", "read", false},   // charlie has no role
	}

	for _, tt := range tests {
		ok, err := svc.Check(tt.Sub, tt.Obj, tt.Act)
		if err != nil {
			t.Errorf("Check failed: %v", err)
		}
		if ok != tt.Expected {
			t.Errorf("Sub %s -> Obj %s (%s): expected %v, got %v", tt.Sub, tt.Obj, tt.Act, tt.Expected, ok)
		}
	}
}

func TestCasbin_InitFail(t *testing.T) {
	// Test error propogation
	_, err := permission.NewService("non_existent_model.conf", "")
	if err == nil {
		t.Error("Expected error for missing model file")
	}
}

func TestCasbin_Management(t *testing.T) {
	// Use a fresh file for management test to avoid polluting the static policy.csv
	// Actually, NewEnforcer works with file paths directly.
	// We'll create a temp policy file.
	
	f, _ := os.CreateTemp("", "policy_*.csv")
	policyPath := f.Name()
	f.Close()
	defer os.Remove(policyPath)

	svc, err := permission.NewService("model.conf", policyPath)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 1. Add Role
	if _, err := svc.AddRoleForUser("david", "manager"); err != nil {
		t.Error(err)
	}

	// 2. Add Policy
	// manager can manage reports
	if _, err := svc.AddPolicy("manager", "/reports/*", "read"); err != nil {
		t.Error(err)
	}

	// 3. Verify
	ok, _ := svc.Check("david", "/reports/2026", "read")
	if !ok {
		t.Error("David should have read access to reports")
	}

	// 4. Remove Policy
	if _, err := svc.RemovePolicy("manager", "/reports/*", "read"); err != nil {
		t.Error(err)
	}
	
	ok, _ = svc.Check("david", "/reports/2026", "read")
	if ok {
		t.Error("David should NOT have read access after removal")
	}
	
	// 5. Get Roles
	roles, _ := svc.GetRolesForUser("david")
	if len(roles) == 0 || roles[0] != "manager" {
		t.Errorf("Expected manager role, got %v", roles)
	}
}
