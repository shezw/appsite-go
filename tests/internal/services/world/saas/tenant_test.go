package saas_test

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"appsite-go/internal/services/world/entity"
	"appsite-go/internal/services/world/saas"
)

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestTenant_CRUD(t *testing.T) {
	db := setupDB(t)
	svc := saas.NewTenantService(db)

	t1 := &entity.Tenant{
		Title:  "Alpha Corp",
		Code:   "alpha",
		Domain: "alpha.com",
		Status: "enabled",
		ExpireAt: time.Now().Add(24 * time.Hour).Unix(),
	}

	// 1. Create
	if err := svc.Create(t1); err != nil {
		t.Fatalf("Failed to create: %v", err)
	}
	if t1.ID == "" {
		t.Error("ID not generated")
	}

	// 2. Get
	got, err := svc.Get(t1.ID)
	if err != nil {
		t.Fatalf("Failed to get: %v", err)
	}
	if got.Title != "Alpha Corp" {
		t.Errorf("Title mismatch")
	}

	// 3. GetByDomain
	byDomain, err := svc.GetByDomain("alpha.com")
	if err != nil {
		t.Errorf("Failed to get by domain: %v", err)
	}
	if byDomain.ID != t1.ID {
		t.Error("Mismatch by domain")
	}

	byCode, err := svc.GetByDomain("alpha")
	if err != nil {
		t.Errorf("Failed to get by code: %v", err)
	}
	if byCode.ID != t1.ID {
		t.Error("Mismatch by code")
	}

	// 4. Update
	updates := map[string]interface{}{"title": "Alpha Inc"}
	if err := svc.Update(t1.ID, updates); err != nil {
		t.Fatalf("Failed to update: %v", err)
	}
	
	got2, _ := svc.Get(t1.ID)
	if got2.Title != "Alpha Inc" {
		t.Error("Update failed")
	}

	// 5. Check Active
	active, err := svc.CheckActive(t1.ID)
	if err != nil || !active {
		t.Errorf("Expected active, got %v, err: %v", active, err)
	}

	// Expired case
	t1.ExpireAt = time.Now().Add(-1 * time.Hour).Unix()
	svc.Update(t1.ID, map[string]interface{}{"expire_at": t1.ExpireAt})
	active, _ = svc.CheckActive(t1.ID)
	if active {
		t.Error("Expected inactive (expired)")
	}

	// Disabled case
	svc.Update(t1.ID, map[string]interface{}{"status": "disabled", "expire_at": time.Now().Add(1*time.Hour).Unix()})
	active, _ = svc.CheckActive(t1.ID)
	if active {
		t.Error("Expected inactive (disabled)")
	}

	// 6. List
	svc.Create(&entity.Tenant{Title: "Beta", Code: "beta"})
	list, total, err := svc.List(1, 10, nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if total < 2 {
		t.Errorf("Expected at least 2 tenants, got %d", total)
	}
	if len(list) < 2 {
		t.Errorf("Expected list len >= 2")
	}
}

func TestTenant_Validation(t *testing.T) {
	db := setupDB(t)
	svc := saas.NewTenantService(db)

	// Missing code
	t2 := &entity.Tenant{Title: "Missing Code"}
	if err := svc.Create(t2); err == nil {
		t.Error("Expected error for missing code")
	}
}
