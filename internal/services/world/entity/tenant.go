// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package entity

import (
	"appsite-go/internal/core/model"
	"appsite-go/pkg/dbs"
)

// Tenant represents a SaaS tenant (or World/Server instance).
type Tenant struct {
	model.Base

	Title    string  `json:"title" gorm:"type:varchar(64);not null;comment:Tenant Name"`
	Domain   string  `json:"domain" gorm:"type:varchar(128);uniqueIndex;comment:Custom Domain"`
	Code     string  `json:"code" gorm:"type:varchar(32);uniqueIndex;comment:Subdomain code"` 
	
	Status   string  `json:"status" gorm:"type:varchar(16);default:'enabled';index"`
	ExpireAt int64   `json:"expire_at" gorm:"index;comment:Subscription Expiry"`
	
	OwnerID  string  `json:"owner_id" gorm:"type:varchar(36);index;comment:Admin User ID"`
	
	Config   dbs.Map `json:"config" gorm:"type:json;comment:Tenant specific settings"`
}

// TableName returns table name
func (Tenant) TableName() string {
	return "sys_tenant"
}
