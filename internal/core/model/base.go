// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base contains common columns for all tables.
type Base struct {
	ID        string `gorm:"primaryKey;type:varchar(32);comment:Unique ID"`
	CreatedAt int64  `gorm:"autoCreateTime;comment:Creation timestamp (seconds)"`
	UpdatedAt int64  `gorm:"autoUpdateTime;comment:Last update timestamp (seconds)"`
}

// BeforeCreate is a GORM hook to generate ID if missing
func (base *Base) BeforeCreate(tx *gorm.DB) error {
	if base.ID == "" {
		// Use UUID, simplified (remove hyphens to fit 32 chars if preferred, or keep 36)
		// Base says varchar(32). UUID is 36.
		// So we must strip hyphens.
		id := uuid.New().String()
		base.ID = strings.ReplaceAll(id, "-", "")
	}
	return nil
}

// SoftDelete adds soft delete capability.
type SoftDelete struct {
	DeletedAt gorm.DeletedAt `gorm:"index;comment:Soft delete timestamp"`
}

// Tenant adds SaaS capability.
type Tenant struct {
	SaasID string `gorm:"index;type:varchar(32);default:'';comment:SaaS Tenant ID"`
}

// TimeFields defines standard time tracking
type TimeFields struct {
	CreatedAt int64 `gorm:"autoCreateTime"`
	UpdatedAt int64 `gorm:"autoUpdateTime"`
}

// ActionHooks defines hooks that entities can implement to intercept actions at the Application layer
// (distinct from GORM's DB-layer hooks).
type ActionHooks interface {
	BeforeAdd(tx *gorm.DB) error
	AfterAdd(tx *gorm.DB) error
	BeforeUpdate(tx *gorm.DB) error
	AfterUpdate(tx *gorm.DB) error
	BeforeDelete(tx *gorm.DB) error
	AfterDelete(tx *gorm.DB) error
}

// Default implementation to allow embedding without implementing all methods? 
// No, interfaces don't work like that. Structs will implement what they need.
// We will use type assertion to check existence.
