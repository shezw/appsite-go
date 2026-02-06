// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package entity

import (
	"appsite-go/internal/core/model"
)

// User represents a registered account
type User struct {
	model.Base
	model.Tenant // SaasID

	// Core Auth
	Username string `gorm:"size:64;uniqueIndex;comment:Login username"`
	Password string `gorm:"size:255;comment:Hashed password"`
	Email    string `gorm:"size:64;uniqueIndex;comment:Login email"`
	Mobile   string `gorm:"size:24;uniqueIndex;comment:Login mobile"`

	// Profile Basic
	Nickname string `gorm:"size:64"`
	Avatar   string `gorm:"size:255"`
	Status   string `gorm:"size:12;default:'enabled'"` // enabled, disabled
	
	// Relations
	GroupID string `gorm:"size:32;default:'100'"`
}

func (User) TableName() string {
	return "user_account" // Matches PHP 'user_account' table
}
