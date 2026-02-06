package contents_test

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"appsite-go/internal/services/contents"
	"appsite-go/internal/services/contents/entity"
)

func setupArticleDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestArticle_CRUD(t *testing.T) {
	db := setupArticleDB(t)
	svc := contents.NewArticleService(db)

	// 1. Create
	article := &entity.Article{
		Title:       "Hello World",
		Description: "First article",
		AuthorID:    "user_001",
		Status:      "enabled",
		Featured:    true,
	}

	if err := svc.Create(article); err != nil {
		t.Fatalf("Failed to create article: %v", err)
	}

	if article.ID == "" {
		t.Fatal("Article ID should be generated")
	}

	// 2. Get
	fetched, err := svc.Get(article.ID)
	if err != nil {
		t.Fatalf("Failed to get article: %v", err)
	}
	if fetched.Title != article.Title {
		t.Errorf("Expected title %s, got %s", article.Title, fetched.Title)
	}
	if !fetched.Featured {
		t.Error("Expected Featured to be true")
	}

	// 3. Update
	updates := map[string]interface{}{
		"title": "Hello Golang",
	}
	if err := svc.Update(article.ID, updates); err != nil {
		t.Fatalf("Failed to update article: %v", err)
	}

	fetchedAfterUpdate, _ := svc.Get(article.ID)
	if fetchedAfterUpdate.Title != "Hello Golang" {
		t.Errorf("Expected title 'Hello Golang', got %s", fetchedAfterUpdate.Title)
	}

	// 4. Increment View
	if err := svc.IncrementView(article.ID); err != nil {
		t.Fatalf("Failed to increment view: %v", err)
	}
	
	fetchedAfterView, _ := svc.Get(article.ID)
	if fetchedAfterView.ViewTimes != 1 {
		t.Errorf("Expected ViewTimes 1, got %d", fetchedAfterView.ViewTimes)
	}

	// 5. List
	svc.Create(&entity.Article{Title: "Article 2", Status: "enabled"})
	svc.Create(&entity.Article{Title: "Article 3", Status: "disabled"})

	listData, total, err := svc.List(1, 10, map[string]interface{}{"status": "enabled"})
	if err != nil {
		t.Fatalf("Failed to list articles: %v", err)
	}

	if total != 2 { // "Hello Golang" and "Article 2" are enabled
		t.Errorf("Expected 2 enabled articles, got %d", total)
	}
	if len(listData) != 2 {
		t.Errorf("Expected length 2, got %d", len(listData))
	}

	// 6. Delete
	if err := svc.Delete(article.ID); err != nil {
		t.Fatalf("Failed to delete article: %v", err)
	}

	_, err = svc.Get(article.ID)
	if err == nil {
		t.Error("Expected error getting deleted article, got nil")
	}
}
