package user_test

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"appsite-go/internal/services/user/entity"
	"appsite-go/internal/services/user/preference"
)

func setupPrefDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	_ = db.AutoMigrate(&entity.UserPreference{})
	return db
}

func TestPreference(t *testing.T) {
	db := setupPrefDB(t)
	svc := preference.NewService(db)

	userID := "user_001"
	key := "theme"

	// 1. Get Non-existent // Should result error from First
	_, err := svc.Get(userID, key)
	if err == nil {
		t.Error("Expected error for non-existent pref")
	}

	// 2. Set (Create)
	err = svc.Set(userID, "", key, "dark")
	if err != nil {
		t.Fatalf("Failed to set pref: %v", err)
	}

	// 3. Get Created
	val, err := svc.Get(userID, key)
	if err != nil {
		t.Fatalf("Failed to get pref: %v", err)
	}
	if val != "dark" {
		t.Errorf("Expected dark, got %s", val)
	}

	// 4. Set (Update)
	err = svc.Set(userID, "", key, "light")
	if err != nil {
		t.Fatalf("Failed to update pref: %v", err)
	}

	val2, _ := svc.Get(userID, key)
	if val2 != "light" {
		t.Errorf("Expected light, got %s", val2)
	}
}
