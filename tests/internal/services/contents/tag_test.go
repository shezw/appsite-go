package contents_test

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"appsite-go/internal/services/contents"
	"appsite-go/internal/services/contents/entity"
)

func setupTagDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestTag_CRUD(t *testing.T) {
	db := setupTagDB(t)
	svc := contents.NewTagService(db)

	// 1. Create
	tag := &entity.Tag{
		Title:  "Golang",
		Type:   "tech",
		Status: "enabled",
	}

	if err := svc.Create(tag); err != nil {
		t.Fatalf("Failed to create tag: %v", err)
	}

	if tag.ID == "" {
		t.Fatal("Tag ID should be generated")
	}

	// 2. Get
	fetched, err := svc.Get(tag.ID)
	if err != nil {
		t.Fatalf("Failed to get tag: %v", err)
	}
	if fetched.Title != tag.Title {
		t.Errorf("Expected title %s, got %s", tag.Title, fetched.Title)
	}

	// 3. Update
	updates := map[string]interface{}{
		"title": "Go Language",
	}
	if err := svc.Update(tag.ID, updates); err != nil {
		t.Fatalf("Failed to update tag: %v", err)
	}

	fetchedAfterUpdate, _ := svc.Get(tag.ID)
	if fetchedAfterUpdate.Title != "Go Language" {
		t.Errorf("Expected title 'Go Language', got %s", fetchedAfterUpdate.Title)
	}

	// 4. List
	svc.Create(&entity.Tag{Title: "Java", Status: "enabled"})
	list, _, err := svc.List(1, 10, nil)
	if err != nil {
		t.Fatalf("Failed to list: %v", err)
	}
	if len(list) < 2 {
		t.Errorf("Expected at least 2 tags, got %d", len(list))
	}

	// 5. Delete
	if err := svc.Delete(tag.ID); err != nil {
		t.Fatalf("Failed to delete tag: %v", err)
	}

	_, err = svc.Get(tag.ID)
	if err == nil {
		t.Error("Expected error getting deleted tag, got nil")
	}
}
