// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package contents_test

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	
	"appsite-go/internal/services/contents"
	"appsite-go/internal/services/contents/entity"
)

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestCategory_Tree(t *testing.T) {
	db := setupDB(t)
	svc := contents.NewCategoryService(db)

	cats := []entity.Category{
		{Title: "Tech", ParentID: ""},        // ID will be generated
		{Title: "Golang", ParentID: ""},      // Will verify Parent later
		{Title: "Python", ParentID: ""},
		{Title: "Frameworks", ParentID: ""},
	}

	// Insert Root "Tech"
	svc.Create(&cats[0])
	techID := cats[0].ID
	
	// Insert Children
	cats[1].ParentID = techID // Go -> Tech
	svc.Create(&cats[1])
	
	cats[2].ParentID = techID // Py -> Tech
	svc.Create(&cats[2])
	
	// Create "Frameworks" under "Golang"
	cats[3].ParentID = cats[1].ID
	svc.Create(&cats[3])

	// Fetch Tree
	nodes, err := svc.Tree("")
	if err != nil {
		t.Fatalf("Tree failed: %v", err)
	}

	if len(nodes) != 1 {
		t.Errorf("Expected 1 root (Tech), got %d", len(nodes))
	}
	if nodes[0].Title != "Tech" {
		t.Errorf("Root title mismatch")
	}
	if len(nodes[0].Children) != 2 {
		t.Errorf("Expected 2 children for Tech, got %d", len(nodes[0].Children))
	}
	
	// Find Golang
	var golang *contents.CategoryNode
	for _, c := range nodes[0].Children {
		if c.Title == "Golang" {
			golang = c
			break
		}
	}
	
	if golang == nil {
		t.Fatal("Golang node not found")
	}
	if len(golang.Children) != 1 {
		t.Errorf("Expected 1 child for Golang, got %d", len(golang.Children))
	}
	if golang.Children[0].Title != "Frameworks" {
		t.Errorf("Grandchild title mismatch")
	}
}

func TestCategory_CRUD(t *testing.T) {
	db := setupDB(t)
	svc := contents.NewCategoryService(db)

	cat := &entity.Category{
		Title: "Test Cat",
		// Code field removed
		Status: "enabled",
	}

	// 1. Create
	if err := svc.Create(cat); err != nil {
		t.Fatalf("Failed to create: %v", err)
	}
	if cat.ID == "" {
		t.Error("ID not generated")
	}

	// 2. Get
	got, err := svc.Get(cat.ID)
	if err != nil {
		t.Fatalf("Failed to get: %v", err)
	}
	if got.Title != "Test Cat" {
		t.Errorf("Title mismatch")
	}

	// 3. Update
	cat.Title = "Updated Cat"
	if err := svc.Update(cat.ID, map[string]interface{}{"title": "Updated Cat"}); err != nil {
		t.Fatalf("Failed to update: %v", err)
	}

	// 4. List
	svc.Create(&entity.Category{Title: "Another Cat", Status: "enabled"})
	list, _, err := svc.List(1, 10, nil)
	if err != nil {
		t.Fatalf("Failed to list: %v", err)
	}
	if len(list) < 2 {
		t.Errorf("Expected at least 2 categories, got %d", len(list))
	}

	// 5. Delete
	if err := svc.Delete(cat.ID); err != nil {
		t.Fatalf("Failed to delete: %v", err)
	}
}
