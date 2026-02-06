package contents_test

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"appsite-go/internal/services/contents"
	"appsite-go/internal/services/contents/entity"
)

func setupCommentDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestComment_CRUD(t *testing.T) {
	db := setupCommentDB(t)
	svc := contents.NewCommentService(db)

	// 1. Create
	comment := &entity.Comment{
		UserID:   "user_123",
		ItemID:   "article_abc",
		ItemType: "article",
		Title:    "Great Article",
		Content:  "I really learned a lot.",
		Status:   "enabled",
	}

	if err := svc.Create(comment); err != nil {
		t.Fatalf("Failed to create comment: %v", err)
	}

	if comment.ID == "" {
		t.Fatal("Comment ID should be generated")
	}

	// 2. Get
	fetched, err := svc.Get(comment.ID)
	if err != nil {
		t.Fatalf("Failed to get comment: %v", err)
	}
	if fetched.Content != comment.Content {
		t.Errorf("Expected content %s, got %s", comment.Content, fetched.Content)
	}

	// 3. Update
	updates := map[string]interface{}{
		"content": "Updated content",
	}
	if err := svc.Update(comment.ID, updates); err != nil {
		t.Fatalf("Failed to update comment: %v", err)
	}

	fetchedAfterUpdate, _ := svc.Get(comment.ID)
	if fetchedAfterUpdate.Content != "Updated content" {
		t.Errorf("Expected content 'Updated content', got %s", fetchedAfterUpdate.Content)
	}

	// 4. List
	svc.Create(&entity.Comment{Content: "Another comment", Status: "approved"})
	list, _, err := svc.List(1, 10, nil)
	if err != nil {
		t.Fatalf("Failed to list: %v", err)
	}
	if len(list) < 2 {
		t.Errorf("Expected at least 2 comments, got %d", len(list))
	}

	// 5. Delete
	if err := svc.Delete(comment.ID); err != nil {
		t.Fatalf("Failed to delete comment: %v", err)
	}

	_, err = svc.Get(comment.ID)
	if err == nil {
		t.Error("Expected error getting deleted comment, got nil")
	}
}
