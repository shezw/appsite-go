package user_test

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"appsite-go/internal/services/user/entity"
	"appsite-go/internal/services/user/info"
)

func setupInfoDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	_ = db.AutoMigrate(&entity.UserInfo{})
	return db
}

func TestProfile(t *testing.T) {
	db := setupInfoDB(t)
	svc := info.NewService(db)

	userID := "user_001"

	// 1. Get Non-existent
	_, err := svc.GetProfile(userID)
	if err == nil {
		t.Error("Expected error for non-existent profile")
	}

	// 2. Update (Create)
	data := map[string]interface{}{
		"real_name": "John Doe",
		"city":      "New York",
	}
	err = svc.UpdateProfile(userID, data)
	if err != nil {
		t.Fatalf("Failed to create profile: %v", err)
	}

	// 3. Get Created
	p, err := svc.GetProfile(userID)
	if err != nil {
		t.Fatalf("Failed to get profile: %v", err)
	}
	if p.RealName != "John Doe" {
		t.Errorf("Expected John Doe, got %s", p.RealName)
	}

	// 4. Update Existing
	data2 := map[string]interface{}{
		"city": "San Francisco",
	}
	err = svc.UpdateProfile(userID, data2)
	if err != nil {
		t.Fatalf("Failed to update profile: %v", err)
	}

	p2, _ := svc.GetProfile(userID)
	if p2.City != "San Francisco" {
		t.Errorf("Expected San Francisco, got %s", p2.City)
	}
	if p2.RealName != "John Doe" {
		t.Errorf("Expected John Doe to remain, got %s", p2.RealName)
	}
}
