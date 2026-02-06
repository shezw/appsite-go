// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package entity

import (
	"appsite-go/internal/core/model"
)

// UserInfo represents extended user details
// Maps to table 'user_info'
type UserInfo struct {
	UserID    string `gorm:"primaryKey;size:32;comment:User ID FK"`
	
	// Identity
	RealName  string `gorm:"size:64;comment:Real Name"`
	IDNumber  string `gorm:"size:32;comment:ID Card Number"`
	
	// Location
	Country   string `gorm:"size:64"`
	Province  string `gorm:"size:64"`
	City      string `gorm:"size:64"`
	Company   string `gorm:"size:128"`
	
	// Social
	WechatID  string `gorm:"size:64;index"`
	WeiboID   string `gorm:"size:64;index"`
	QQID      string `gorm:"size:64;index"`
	AppleUUID string `gorm:"size:64;index"`
	
	// Device
	DeviceID  string `gorm:"size:64;index"`
	
	// Advanced
	VIP       int   `gorm:"default:0"`
	VIPExpire int64 `gorm:"default:0"`
	Gallery   string `gorm:"type:text;comment:JSON array of image URLs"` // JSON
}

func (UserInfo) TableName() string {
	return "user_info"
}

// UserGroup represents user roles/groups
// Maps to table 'user_group'
type UserGroup struct {
	model.Base
	
	ParentID    string `gorm:"size:32;default:'0';index"`
	Type        string `gorm:"size:32"`
	Level       int    `gorm:"default:0"`
	
	GroupName   string `gorm:"size:64"`
	Description string `gorm:"size:255"`
	MenuAccess  string `gorm:"type:text;comment:JSON of accessible menu IDs"`
	
	Status      string `gorm:"size:12;default:'enabled'"`
	Sort        int    `gorm:"default:0"`
}

func (UserGroup) TableName() string {
	return "user_group"
}

// UserPreference represents key-value user settings
// Maps to table 'user_preference'
type UserPreference struct {
	model.Base
	
	UserID    string `gorm:"size:32;index"`
	SaasID    string `gorm:"size:32;index"`
	
	KeyID     string `gorm:"size:64;index"`
	Content   string `gorm:"type:text"`
	Desc      string `gorm:"size:255"`
	Version   string `gorm:"size:32"`
	
	Status    string `gorm:"size:12;default:'enabled'"`
}

func (UserPreference) TableName() string {
	return "user_preference"
}
