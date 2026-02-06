package relation_test

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"appsite-go/internal/services/relation"
)

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestRelation(t *testing.T) {
	db := setupDB(t)
	svc := relation.NewService(db)

	userA := "user_a"
	userB := "user_b"

	// 1. Bind (A follows B)
	err := svc.Bind(userA, "user", userB, "user", "follow")
	if err != nil {
		t.Fatalf("Failed to bind: %v", err)
	}

	// 2. Check
	if !svc.Check(userA, "user", userB, "user", "follow") {
		t.Error("Expected relation to exist")
	}

	// 3. Count (B's followers)
	count := svc.Count(userB, "user", "follow")
	if count != 1 {
		t.Errorf("Expected 1 follower, got %d", count)
	}

	// 4. List (A's following)
	list, _, err := svc.ListRelations(userA, "user", "follow", 1, 10)
	if err != nil {
		t.Fatalf("Failed to list: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("Expected 1 relation in list, got %d", len(list))
	}
	if list[0].RelationID != userB {
		t.Errorf("Expected user_b, got %s", list[0].RelationID)
	}

	// 5. Unbind
	err = svc.Unbind(userA, "user", userB, "user", "follow")
	if err != nil {
		t.Fatalf("Failed to unbind: %v", err)
	}

	if svc.Check(userA, "user", userB, "user", "follow") {
		t.Error("Expected relation to be removed")
	}
}
