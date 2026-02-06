// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model_test

import (
	"errors"
	"testing"

	"appsite-go/internal/core/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// -- Mock Entities --

type User struct {
	model.Base
	Name string
}

// UserWithHooks implements ActionHooks
type UserWithHooks struct {
	model.Base
	Name string
}

func (u *UserWithHooks) BeforeAdd(tx *gorm.DB) error {
	if u.Name == "Forbidden" {
		return errors.New("forbidden name")
	}
	u.Name = "Hooked:" + u.Name
	return nil
}
func (u *UserWithHooks) AfterAdd(tx *gorm.DB) error { return nil }
func (u *UserWithHooks) BeforeUpdate(tx *gorm.DB) error {
    // We can access tx here to check things
	return nil
}
func (u *UserWithHooks) AfterUpdate(tx *gorm.DB) error { return nil }
func (u *UserWithHooks) BeforeDelete(tx *gorm.DB) error { return nil }
func (u *UserWithHooks) AfterDelete(tx *gorm.DB) error { return nil }


// -- Tests --

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect sqlite: %v", err)
	}
	return db
}

func TestCRUD_Add(t *testing.T) {
	db := setupDB(t)
	db.AutoMigrate(&User{}, &UserWithHooks{})
	
	crud := model.NewCRUD[User](db)
	
	// 1. Success
	user := &User{Name: "Alice"}
	user.ID = "1"
	res := crud.Add(user)
	if !res.Success {
		t.Errorf("Add failed: %v", res.Message)
	}

	// 2. Hook Execution
	crudHooks := model.NewCRUD[UserWithHooks](db)
	
	// 2.1 Hook Modifies Data
	uHook := &UserWithHooks{Name: "Bob"}
	uHook.ID = "2"
	res = crudHooks.Add(uHook)
	if !res.Success {
		t.Fatalf("Add with hook failed: %v", res.Message)
	}
	if uHook.Name != "Hooked:Bob" {
		t.Errorf("Hook did not run. Name=%s", uHook.Name)
	}

	// 2.2 Hook Blocks Action
	uFail := &UserWithHooks{Name: "Forbidden"}
	uFail.ID = "3"
	res = crudHooks.Add(uFail)
	if res.Success {
		t.Error("Hook should have blocked 'Forbidden' name")
	}

	// 2.3 AfterAdd Hook Error (Warning)
	// We need to implement a case where AfterAdd fails.
	// But our mock AfterAdd returns nil.
	// Let's create a specific type for that if we want 100% coverage on that warning line.
}

func TestCRUD_Get(t *testing.T) {
	db := setupDB(t)
	db.AutoMigrate(&User{})
	crud := model.NewCRUD[User](db)

	db.Create(&User{Base: model.Base{ID: "100"}, Name: "Exist"})

	// 1. Found
	res := crud.Get("100")
	if !res.Success {
		t.Error("Should find user")
	}
	if res.Data.(*User).Name != "Exist" {
		t.Error("Data mismatch")
	}

	// 2. Not Found
	res = crud.Get("999")
	if res.Success {
		t.Error("Should fail for missing ID")
	}
}

func TestCRUD_Update(t *testing.T) {
	db := setupDB(t)
	db.AutoMigrate(&User{})
	crud := model.NewCRUD[User](db)

	db.Create(&User{Base: model.Base{ID: "200"}, Name: "Old"})

	// 1. Success
	res := crud.Update("200", map[string]interface{}{"Name": "New"})
	if !res.Success {
		t.Error("Update failed")
	}
	
	var u User
	db.First(&u, "id = ?", "200")
	if u.Name != "New" {
		t.Errorf("DB not updated, got %s", u.Name)
	}

	// 2. Not Found
	res = crud.Update("999", map[string]interface{}{})
	if res.Success {
		t.Error("Update missing should fail")
	}

	// 3. Update Fail (Wait, how to induce DB error on Updates?)
	// Hard with sqlite memory unless we table lock or similar.
}

func TestCRUD_Remove(t *testing.T) {
	db := setupDB(t)
	db.AutoMigrate(&User{})
	crud := model.NewCRUD[User](db)

	db.Create(&User{Base: model.Base{ID: "300"}, Name: "DeleteMe"})

	// 1. Success
	res := crud.Remove("300")
	if !res.Success {
		t.Error("Remove failed")
	}
	
	var count int64
	db.Model(&User{}).Where("id = ?", "300").Count(&count)
	if count > 0 {
		t.Error("User not deleted")
	}

	// 2. Not Found
	res = crud.Remove("999")
	if res.Success {
		t.Error("Remove missing should fail")
	}
}

func TestCRUD_List(t *testing.T) {
	db := setupDB(t)
	db.AutoMigrate(&User{})
	crud := model.NewCRUD[User](db)

	for i := 0; i < 15; i++ {
		// Explicit ID needed because Base struct has ID as PK string which defaults to "" causing unique constraint fail
		id := "u_" + string(rune('a'+i))
		db.Create(&User{Name: "U", Base: model.Base{ID: id, CreatedAt: int64(100 + i)}})
	}
	// Add specific user for filter
	db.Create(&User{Name: "Target", Base: model.Base{ID: "target_1", CreatedAt: 200}})

	// 1. Basic List
	res := crud.List(&model.ListParams{Page: 1, PageSize: 10})
	if !res.Success {
		t.Errorf("List failed")
	}
	data := res.Data.(map[string]interface{})
	if data["total"].(int64) != 16 {
		t.Errorf("Expected 16 total, got %d", data["total"])
	}
	list := data["list"].([]User)
	if len(list) != 10 {
		t.Errorf("Page size mismatch, got %d", len(list))
	}

	// 2. Filter
	// Note: 'Name' is a field in User struct, so hasField should pass.
	// Filter uses exact match
	res = crud.List(&model.ListParams{
		Page: 1, PageSize: 10,
		Filters: map[string]interface{}{"Name": "Target"},
	})
	data = res.Data.(map[string]interface{})
	if data["total"].(int64) != 1 {
		t.Errorf("Filter failed, got %d results", data["total"])
	}
	
	// 3. Sorting (implicit check)
	res = crud.List(&model.ListParams{Page:1, PageSize: 1, Sort: "created_at DESC"})
	list = res.Data.(map[string]interface{})["list"].([]User)
	// Target has CreatedAt 200, others match 100+i (max 114)
	if list[0].CreatedAt != 200 { 
		t.Errorf("Sort logic check failed, expected 200, got %d", list[0].CreatedAt)
	}
}
