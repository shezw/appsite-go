package contents_test

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"appsite-go/internal/services/contents"
	"appsite-go/internal/services/contents/entity"
)

func setupBannerDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestBanner_CRUD(t *testing.T) {
	db := setupBannerDB(t)
	svc := contents.NewBannerService(db)

	// 1. Create
	banner := &entity.Banner{
		Title:    "Home Banner",
		Position: "home_top",
		Status:   "enabled",
	}

	if err := svc.Create(banner); err != nil {
		t.Fatalf("Failed to create banner: %v", err)
	}

	if banner.ID == "" {
		t.Fatal("Banner ID should be generated")
	}

	// 2. Get
	fetched, err := svc.Get(banner.ID)
	if err != nil {
		t.Fatalf("Failed to get banner: %v", err)
	}
	if fetched.Title != banner.Title {
		t.Errorf("Expected title %s, got %s", banner.Title, fetched.Title)
	}

	// 3. Update
	updates := map[string]interface{}{
		"title": "Home Banner Updated",
	}
	if err := svc.Update(banner.ID, updates); err != nil {
		t.Fatalf("Failed to update banner: %v", err)
	}

	fetchedAfterUpdate, _ := svc.Get(banner.ID)
	if fetchedAfterUpdate.Title != "Home Banner Updated" {
		t.Errorf("Expected title 'Home Banner Updated', got %s", fetchedAfterUpdate.Title)
	}

	// 4. Increment Click & View
	if err := svc.IncrementClick(banner.ID); err != nil {
		t.Fatalf("Failed to increment click: %v", err)
	}
	if err := svc.IncrementView(banner.ID); err != nil {
		t.Fatalf("Failed to increment view: %v", err)
	}
	
	fetchedAfterClick, _ := svc.Get(banner.ID)
	if fetchedAfterClick.ClickTimes != 1 {
		t.Errorf("Expected ClickTimes 1, got %d", fetchedAfterClick.ClickTimes)
	}
	if fetchedAfterClick.ViewTimes != 1 {
		t.Errorf("Expected ViewTimes 1, got %d", fetchedAfterClick.ViewTimes)
	}

	// 5. List
	svc.Create(&entity.Banner{Title: "Banner 2", Status: "enabled"})
	list, _, err := svc.List(1, 10, map[string]interface{}{"status": "enabled"})
	if err != nil {
		t.Fatalf("Failed to list: %v", err)
	}
	if len(list) < 2 {
		t.Errorf("Expected at least 2 banners, got %d", len(list))
	}

	// 6. Delete
	if err := svc.Delete(banner.ID); err != nil {
		t.Fatalf("Failed to delete banner: %v", err)
	}

	_, err = svc.Get(banner.ID)
	if err == nil {
		t.Error("Expected error getting deleted banner, got nil")
	}
}
