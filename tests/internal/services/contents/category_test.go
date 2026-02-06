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
