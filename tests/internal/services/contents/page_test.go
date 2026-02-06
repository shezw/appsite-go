package contents_test

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"appsite-go/internal/services/contents"
	"appsite-go/internal/services/contents/entity"
)

func setupPageDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestPage_CRUD(t *testing.T) {
	db := setupPageDB(t)
	svc := contents.NewPageService(db)

	// 1. Create
	page := &entity.Page{
		Title:  "About Us",
		Alias:  "about-us",
		Status: "enabled",
	}

	if err := svc.Create(page); err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	if page.ID == "" {
		t.Fatal("Page ID should be generated")
	}

	// 2. Get by alias
	fetched, err := svc.GetByAlias("about-us")
	if err != nil {
		t.Fatalf("Failed to get page by alias: %v", err)
	}
	if fetched.ID != page.ID {
		t.Errorf("Expected ID %s, got %s", page.ID, fetched.ID)
	}

	// 3. Update
	updates := map[string]interface{}{
		"title": "About Us Updated",
	}
	if err := svc.Update(page.ID, updates); err != nil {
		t.Fatalf("Failed to update page: %v", err)
	}

	fetchedAfterUpdate, _ := svc.Get(page.ID)
	if fetchedAfterUpdate.Title != "About Us Updated" {
		t.Errorf("Expected title 'About Us Updated', got %s", fetchedAfterUpdate.Title)
	}

	// 4. Increment View
	if err := svc.IncrementView(page.ID); err != nil {
		t.Fatalf("Failed to increment view: %v", err)
	}
	
	fetchedAfterView, _ := svc.Get(page.ID)
	if fetchedAfterView.ViewTimes != 1 {
		t.Errorf("Expected ViewTimes 1, got %d", fetchedAfterView.ViewTimes)
	}

	// 5. Delete
	if err := svc.Delete(page.ID); err != nil {
		t.Fatalf("Failed to delete page: %v", err)
	}

	_, err = svc.Get(page.ID)
	if err == nil {
		t.Error("Expected error getting deleted page, got nil")
	}
}
