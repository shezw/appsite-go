package shieldword_test

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"appsite-go/internal/services/shieldword"
	"appsite-go/internal/services/shieldword/entity"
)

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestFilter(t *testing.T) {
	db := setupDB(t)
	svc := shieldword.NewService(db)

	// 1. Add Words
	badWords := []string{"bad", "evil", "ugly"}
	for _, w := range badWords {
		svc.Create(&entity.Word{Title: w, Status: "enabled"})
	}

	// 2. Check
	if !svc.Check("This is a bad day") {
		t.Error("Expected to detect 'bad'")
	}
	if svc.Check("This is a good day") {
		t.Error("Did not expect to detect anything in good content")
	}

	// 3. Replace
	result := svc.Replace("This is a bad and evil day")
	expected := "This is a *** and **** day"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	// 4. Cache Invalidation
	// Create another
	word := &entity.Word{Title: "specialword", Status: "enabled"}
	svc.Create(word)
	if !svc.Check("This is specialword") {
		t.Error("Expected to detect newly added word")
	}

	// Delete
	if err := svc.Delete(word.ID); err != nil {
		t.Fatalf("Failed to delete word: %v", err)
	}

	// Verify DB state
	var count int64
	db.Model(&entity.Word{}).Where("id = ?", word.ID).Count(&count)
	if count > 0 {
		t.Errorf("Word still exists in DB after delete. Count: %d", count)
	}

	if svc.Check("This is specialword") {
		t.Error("Did not expect to detect deleted word")
	}
}
