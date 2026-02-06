package form_test

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"appsite-go/internal/services/form"
	"appsite-go/internal/services/form/entity"
	"appsite-go/pkg/dbs"
)

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestSubmission(t *testing.T) {
	db := setupDB(t)
	svc := form.NewSubmissionService(db)

	formData := make(dbs.Map)
	formData["name"] = "John Doe"
	formData["reason"] = "I want to apply"

	req := &entity.Request{
		UserID: "user_01",
		Form:   formData,
	}

	// 1. Submit
	err := svc.Submit(req)
	if err != nil {
		t.Fatalf("Failed to submit: %v", err)
	}

	if req.ID == "" {
		t.Fatal("Request ID should be generated")
	}
	if req.Status != "pending" {
		t.Errorf("Expected status pending, got %s", req.Status)
	}

	// 2. Get
	fetched, err := svc.Get(req.ID)
	if err != nil {
		t.Fatalf("Failed to get request: %v", err)
	}
	if fetched.Form["name"] != "John Doe" {
		t.Errorf("Expected Form name John Doe, got %v", fetched.Form["name"])
	}

	// 3. Review (Approve)
	err = svc.Review(req.ID, true)
	if err != nil {
		t.Fatalf("Failed to review: %v", err)
	}

	fetchedReviewed, _ := svc.Get(req.ID)
	if fetchedReviewed.Status != "applied" {
		t.Errorf("Expected status applied, got %s", fetchedReviewed.Status)
	}

	// 4. List
	list, total, err := svc.List(1, 10, map[string]interface{}{"status": "applied"})
	if err != nil {
		t.Fatalf("Failed to list: %v", err)
	}

	if total != 1 {
		t.Errorf("Expected total 1, got %d", total)
	}
	if len(list) != 1 {
		t.Errorf("Expected list length 1, got %d", len(list))
	}
}
