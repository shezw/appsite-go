// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package entity

import "appsite-go/internal/core/model"

// Category represents a taxonomy item
type Category struct {
	model.Base
	model.Tenant

	Title       string `gorm:"size:64;index;comment:Category Name"`
	Alias       string `gorm:"size:24;index;comment:Unique Alias/Slug"`
	AuthorID    string `gorm:"size:32;index"`
	ParentID    string `gorm:"size:32;index;default:''"`
	Type        string `gorm:"size:32;index;comment:Category Type (e.g. article, video)"`
	
	Description string `gorm:"size:255"`
	Cover       string `gorm:"size:255"` // Image URL

	Sort        int    `gorm:"default:0;index"`
	Featured    bool   `gorm:"default:false;index"`
	Status      string `gorm:"size:12;default:'enabled'"`
}

func (Category) TableName() string {
	return "item_category"
}
