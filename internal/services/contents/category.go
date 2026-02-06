// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package contents

import (
	"gorm.io/gorm"
	
	"appsite-go/internal/core/model"
	"appsite-go/internal/services/contents/entity"
)

// CategoryService handles taxonomy operations
type CategoryService struct {
	db   *gorm.DB
	repo *model.CRUD[entity.Category]
}

// NewCategoryService initializes the service
func NewCategoryService(db *gorm.DB) *CategoryService {
	if db != nil {
		_ = db.AutoMigrate(&entity.Category{})
	}
	return &CategoryService{
		db:   db,
		repo: model.NewCRUD[entity.Category](db),
	}
}

// Create adds a new category
func (s *CategoryService) Create(cat *entity.Category) error {
	res := s.repo.Add(cat)
	return res.Error
}

// Update modifies an existing category
func (s *CategoryService) Update(id string, updates map[string]interface{}) error {
	res := s.repo.Update(id, updates)
	return res.Error
}

// Delete removes a category
func (s *CategoryService) Delete(id string) error {
	res := s.repo.Remove(id)
	return res.Error
}

// Get retrieves a category by ID
func (s *CategoryService) Get(id string) (*entity.Category, error) {
	res := s.repo.Get(id)
	if !res.Success {
		return nil, res.Error
	}
	return res.Data.(*entity.Category), nil
}

// List returns categories with filters
func (s *CategoryService) List(page, size int, filters map[string]interface{}) ([]entity.Category, int64, error) {
	res := s.repo.List(&model.ListParams{
		Page:     page,
		PageSize: size,
		Filters:  filters,
		Sort:     "sort desc, created_at desc",
	})
	
	if !res.Success {
		return nil, 0, res.Error
	}
	
	data := res.Data.(map[string]interface{})
	return data["list"].([]entity.Category), data["total"].(int64), nil
}

// Tree returns a hierarchical structure of categories
// Note: This is a heavy operation, typically we run this for small datasets or cache it.
type CategoryNode struct {
	entity.Category
	Children []*CategoryNode `json:"children"`
}

func (s *CategoryService) Tree(rootID string) ([]*CategoryNode, error) {
	// Fetch all enabled categories
	// Optimization: If tree is too big, fetch only levels?
	// For now, fetch all.
	var all []entity.Category
	err := s.db.Where("status = ?", "enabled").Order("sort desc").Find(&all).Error
	if err != nil {
		return nil, err
	}
	
	// Build map
	nodeMap := make(map[string]*CategoryNode)
	var roots []*CategoryNode
	
	// Pass 1: Create nodes
	for i := range all {
		// Avoid implicit memory aliasing in loops (Go < 1.22)
		cat := all[i] 
		nodeMap[cat.ID] = &CategoryNode{Category: cat}
	}
	
	// Pass 2: Link
	for _, cat := range all {
		node := nodeMap[cat.ID]
		if cat.ParentID == rootID || (rootID == "" && cat.ParentID == "") || (rootID == "0" && cat.ParentID == "0") {
			roots = append(roots, node)
		} else {
			if parent, ok := nodeMap[cat.ParentID]; ok {
				parent.Children = append(parent.Children, node)
			} else if rootID == "" {
				// Parent not found or disabled? Treat as root?
				// Typically orphaned nodes. Skip or add to root.
				// Let's add to roots for visibility if strict mode is off.
				// But strict tree logic says: if parent missing in this set, it's a root of this partial tree?
				// Simpler: Just handle normal parent logic.
			}
		}
	}
	
	return roots, nil
}
