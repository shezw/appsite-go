// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model_test

import (
	"fmt"
	"testing"

	"appsite-go/internal/core/model"
	"appsite-go/pkg/utils/orm"
    "appsite-go/internal/core/setting"

	"gorm.io/gorm"
)

// -- MySQL Specific Models --

type Author struct {
	model.Base
	Name  string
	Posts []Post `gorm:"foreignKey:AuthorID"`
}

type Post struct {
	model.Base
	Title     string
	AuthorID  string
	Content   string
}

func setupMySQL(t *testing.T) *gorm.DB {
	cfg := &setting.DatabaseConfig{
		User:     "testuser",
		Password: "testpass",
		Host:     "127.0.0.1",
		Port:     "3306",
		Name:     "testdb",
		Charset:  "utf8mb4",
        Type:     "mysql",
	}

	db, err := orm.NewMySQLConnection(cfg)
	if err != nil {
		t.Skipf("Skipping MySQL tests: %v", err)
	}
    
    // Clean up tables
    db.Migrator().DropTable(&Author{}, &Post{})
    err = db.AutoMigrate(&Author{}, &Post{})
    if err != nil {
        t.Fatalf("Failed to migrate mysql tables: %v", err)
    }
    
	return db
}

func TestMySQL_CRUD_Join(t *testing.T) {
	db := setupMySQL(t)

    // 1. Basic CRUD (Author)
    authorRepo := model.NewCRUD[Author](db)
    
    author := &Author{Name: "John Doe"}
    author.ID = "auth_1"
    
    // Create
    if res := authorRepo.Add(author); !res.Success {
        t.Fatalf("MySQL Add Author failed: %s", res.Message)
    }
    
    // Read
    if res := authorRepo.Get("auth_1"); !res.Success {
        t.Fatalf("MySQL Get Author failed")
    } else if res.Data.(*Author).Name != "John Doe" {
        t.Errorf("Name mismatch")
    }

    // Update
    if res := authorRepo.Update("auth_1", map[string]interface{}{"Name": "John Smith"}); !res.Success {
        t.Errorf("MySQL Update failed")
    }

    // 2. Multi-Row (Posts)
    postRepo := model.NewCRUD[Post](db)
    
    // Create the second author to satisfy FK
    authorRepo.Add(&Author{Base: model.Base{ID: "other_2"}, Name: "Jane Doe"})

    posts := []*Post{
        {Base: model.Base{ID: "p1"}, Title: "Go Lang", AuthorID: "auth_1"},
        {Base: model.Base{ID: "p2"}, Title: "MySQL", AuthorID: "auth_1"},
        {Base: model.Base{ID: "p3"}, Title: "Other", AuthorID: "other_2"},
    }
    
    for _, p := range posts {
        if res := postRepo.Add(p); !res.Success {
            t.Errorf("Add Post %s failed: %v", p.Title, res.Error)
        }
    }

    // 3. Filtering & Sorting
    listRes := postRepo.List(&model.ListParams{
        Filters: map[string]interface{}{"author_id": "auth_1"}, // Use column name "author_id"
    })
    
    if !listRes.Success {
        t.Fatalf("List Posts failed")
    }
    
    data := listRes.Data.(map[string]interface{})
    list := data["list"].([]Post)
    
    if data["total"].(int64) != 2 {
        t.Errorf("Filter total should be 2, got %d", data["total"])
    }
    if list[0].Title != "Go Lang" { // "Go Lang" < "MySQL"
        t.Errorf("Sort failed, expected Go Lang first")
    }

    // 4. Joint Query
    // We want to fetch Authors including their Posts.
    // GORM uses Preload for this usually.
    // Or we strictly use Joins for specific fields.
    // Let's test Preload via custom query logic or extended CRUD.
    
    // Since Generic CRUD.List is simple, let's verify we can use the DB instance from CRUD to do complex join.
    // Or if we strictly follow "Join Query" logic:
    
    var results []struct {
        AuthorName string
        PostTitle  string
    }
    
    err := authorRepo.DB.Table("author").
        Select("author.name as author_name, post.title as post_title").
        Joins("left join post on post.author_id = author.id").
        Where("author.id = ?", "auth_1").
        Scan(&results).Error
        
    if err != nil {
        t.Fatalf("Join query failed: %v", err)
    }
    
    if len(results) != 2 {
        t.Errorf("Expected 2 joined rows, got %d", len(results))
    }
    
    // 5. Delete
    // Clean up children first to satisfy FK constraint
    postRepo.Remove("p1")
    postRepo.Remove("p2")

    if res := authorRepo.Remove("auth_1"); !res.Success {
        t.Errorf("Remove author failed")
    }
    // Verify Cascade? GORM default might not cascade unless configured in foreign key.
    // Here we just test the manual execution.
}

func TestMySQL_Transaction(t *testing.T) {
    db := setupMySQL(t)
    
    // Transaction Test
    err := db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(&Author{Base: model.Base{ID: "tx_1"}, Name: "TX User"}).Error; err != nil {
            return err
        }
        // Force error to rollback
        return fmt.Errorf("rollback")
    })
    
    if err == nil {
        t.Error("Transaction should have returned error")
    }
    
    // Verify rollback
    var count int64
    db.Model(&Author{}).Where("id = ?", "tx_1").Count(&count)
    if count != 0 {
        t.Error("Transaction rollback failed, user exists")
    }
}
