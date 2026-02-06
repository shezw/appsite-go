package contents_test

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"appsite-go/internal/services/contents"
	"appsite-go/internal/services/contents/entity"
)

func setupMediaDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestMedia_CRUD(t *testing.T) {
	db := setupMediaDB(t)
	svc := contents.NewMediaService(db)

	// 1. Create
	media := &entity.Media{
		Type:   "image",
		URL:    "https://example.com/image.jpg",
		Size:   1024,
		Status: "enabled",
	}

	if err := svc.Create(media); err != nil {
		t.Fatalf("Failed to create media: %v", err)
	}

	if media.ID == "" {
		t.Fatal("Media ID should be generated")
	}

	// 2. Get
	fetched, err := svc.Get(media.ID)
	if err != nil {
		t.Fatalf("Failed to get media: %v", err)
	}
	if fetched.URL != media.URL {
		t.Errorf("Expected URL %s, got %s", media.URL, fetched.URL)
	}

	// 3. Update
	updates := map[string]interface{}{
		"size": 2048,
	}
	if err := svc.Update(media.ID, updates); err != nil {
		t.Fatalf("Failed to update media: %v", err)
	}

	fetchedAfterUpdate, _ := svc.Get(media.ID)
	if fetchedAfterUpdate.Size != 2048 {
		t.Errorf("Expected size 2048, got %d", fetchedAfterUpdate.Size)
	}

	// 4. List
	svc.Create(&entity.Media{Type: "video", URL: "v.mp4", Status: "enabled"})
	list, _, err := svc.List(1, 10, nil)
	if err != nil {
		t.Fatalf("Failed to list: %v", err)
	}
	if len(list) < 2 {
		t.Errorf("Expected at least 2 media, got %d", len(list))
	}

	// 5. Delete
	if err := svc.Delete(media.ID); err != nil {
		t.Fatalf("Failed to delete media: %v", err)
	}

	_, err = svc.Get(media.ID)
	if err == nil {
		t.Error("Expected error getting deleted media, got nil")
	}
}
