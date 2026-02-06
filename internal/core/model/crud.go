// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"errors"

	"gorm.io/gorm"
	"appsite-go/pkg/utils/orm"
)

// CRUD implements the standard operations for any entity T.
// T should be a struct pointer type usually, but generics constraints handles T as the struct type.
// We use *T in arguments.
type CRUD[T any] struct {
	DB *gorm.DB
}

// NewCRUD creates a new operator
func NewCRUD[T any](db *gorm.DB) *CRUD[T] {
	return &CRUD[T]{DB: db}
}

// Add inserts a new entity
func (c *CRUD[T]) Add(entity *T) *Result {
	// Application Layer Hooks
	if hook, ok := any(entity).(interface{ BeforeAdd(*gorm.DB) error }); ok {
		if err := hook.BeforeAdd(c.DB); err != nil {
			return &Result{Success: false, Error: err, Message: err.Error()}
		}
	}

	// DB Operation
	if err := c.DB.Create(entity).Error; err != nil {
		return &Result{Success: false, Error: err, Message: "Database Create Failed"}
	}

	// After Hook
	if hook, ok := any(entity).(interface{ AfterAdd(*gorm.DB) error }); ok {
		if err := hook.AfterAdd(c.DB); err != nil {
			// Note: Entity is already saved. Error here is non-blocking for data, but reported.
			return &Result{Success: true, Data: entity, Message: "Created with warning: " + err.Error()}
		}
	}

	return &Result{Success: true, Data: entity}
}

// Get retrieves an entity by ID
func (c *CRUD[T]) Get(id string) *Result {
	var entity T
	if err := c.DB.First(&entity, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &Result{Success: false, Error: err, Message: "Not Found"}
		}
		return &Result{Success: false, Error: err, Message: "Database Error"}
	}
	return &Result{Success: true, Data: &entity}
}

// Update updates fields of an entity by ID
func (c *CRUD[T]) Update(id string, updates map[string]interface{}) *Result {
	var entity T
	// Check existence
	if err := c.DB.First(&entity, "id = ?", id).Error; err != nil {
		return &Result{Success: false, Error: err, Message: "Not Found"}
	}

	// Before Hook logic requires an instance.
	// Since updates are a map, we can't easily call "BeforeUpdate" on the *struct* with the new values 
	// unless we marshal them into the struct. 
	// However, usually we might run logic on the *existing* entity or the *updates*.
	// Simpler approach: If the ENTITY type supports BeforeUpdate, we call it on the FETCHED entity
	// passing the DB transaction scope which might contain updates.
	
	if hook, ok := any(&entity).(interface{ BeforeUpdate(*gorm.DB) error }); ok {
		if err := hook.BeforeUpdate(c.DB); err != nil {
			return &Result{Success: false, Error: err, Message: err.Error()}
		}
	}

	if err := c.DB.Model(&entity).Updates(updates).Error; err != nil {
		return &Result{Success: false, Error: err, Message: "Update Failed"}
	}

	return &Result{Success: true, Data: &entity}
}

// Remove deletes an entity
func (c *CRUD[T]) Remove(id string) *Result {
	var entity T
	// We need to find it first to trigger GORM hooks properly if they exist on the model
	if err := c.DB.First(&entity, "id = ?", id).Error; err != nil {
		return &Result{Success: false, Error: err, Message: "Not Found"}
	}

	if hook, ok := any(&entity).(interface{ BeforeDelete(*gorm.DB) error }); ok {
		if err := hook.BeforeDelete(c.DB); err != nil {
			return &Result{Success: false, Error: err, Message: err.Error()}
		}
	}

	if err := c.DB.Delete(&entity).Error; err != nil {
		return &Result{Success: false, Error: err, Message: "Delete Failed"}
	}

	return &Result{Success: true, Data: id}
}

// List retrieves a paginated list
func (c *CRUD[T]) List(params *ListParams) *Result {
	var items []T
	var total int64
	
	db := c.DB.Model(new(T))

	// Apply Filters
	if params.Filters != nil && len(params.Filters) > 0 {
		// Pass map directly to GORM. 
		// keys must be column names or GORM handles them if they match columns.
		db = db.Where(params.Filters)
	}

	// Count
	db.Count(&total)

	// Scopes
	db = db.Scopes(orm.Paginate(params.Page, params.PageSize))

	if params.Sort != "" {
		db = db.Order(params.Sort)
	} else {
		db = db.Order("created_at DESC")
	}

	if err := db.Find(&items).Error; err != nil {
		return &Result{Success: false, Error: err, Message: "List Failed"}
	}

	return &Result{Success: true, Data: map[string]interface{}{
		"list":  items,
		"total": total,
		"page":  params.Page,
		"size":  params.PageSize,
	}}
}
