package user_test

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"appsite-go/internal/services/user/entity"
	"appsite-go/internal/services/user/group"
)

func setupGroupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	_ = db.AutoMigrate(&entity.UserGroup{})
	return db
}

func TestGroup(t *testing.T) {
	db := setupGroupDB(t)
	svc := group.NewService(db)

	// 1. Create
	g, err := svc.Create("Admin", "Super admin", 99)
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}
	if g.ID == "" {
		t.Error("Expected generated ID")
	}

	// 2. List
	list, err := svc.List()
	if err != nil {
		t.Fatalf("Failed to list: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("Expected 1 group, got %d", len(list))
	}
}
