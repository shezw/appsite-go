// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package entity

import (
	"appsite-go/internal/core/model"
	"appsite-go/pkg/dbs"
)

// Product represents a Standard Product Unit (SPU)
type Product struct {
	model.Base
	model.Tenant

	Title      string  `json:"title" gorm:"type:varchar(128);not null;index"`
	SubTitle   string  `json:"sub_title" gorm:"type:varchar(255)"`
	CategoryID string  `json:"category_id" gorm:"type:varchar(36);index"`
	BrandID    string  `json:"brand_id" gorm:"type:varchar(36);index"`
	
	Cover      string  `json:"cover" gorm:"type:varchar(255)"`
	Images     dbs.Map `json:"images" gorm:"type:json"` // Array of strings
	
	Price      int64   `json:"price" gorm:"comment:Display Price in cents"`
	MarketPrice int64  `json:"market_price" gorm:"comment:Original Price"`
	
	Status     string  `json:"status" gorm:"type:varchar(16);default:'offline';index"` // on_sale, offline, sold_out
	
	Content    string  `json:"content" gorm:"type:text"`
	Specs      dbs.Map `json:"specs" gorm:"type:json;comment:Specification Template"` // e.g. [{"name":"Color", "values":["Red","Blue"]}]
}

// TableName returns table name
func (Product) TableName() string {
	return "shop_product"
}

// SKU represents a Stock Keeping Unit
type SKU struct {
	model.Base
	model.Tenant

	ProductID string  `json:"product_id" gorm:"type:varchar(36);index;not null"`
	Code      string  `json:"code" gorm:"type:varchar(64);index;comment:Unique SKU Code"`
	
	Title     string  `json:"title" gorm:"type:varchar(128);comment:Specific Name"`
	Cover     string  `json:"cover" gorm:"type:varchar(255)"`
	
	Price     int64   `json:"price" gorm:"not null"`
	CostPrice int64   `json:"cost_price"`
	Stock     int     `json:"stock" gorm:"default:0"`
	
	Specs     dbs.Map `json:"specs" gorm:"type:json;comment:Specific Spec Values"` // e.g. {"Color":"Red", "Size":"XL"}
	
	Status    string  `json:"status" gorm:"type:varchar(16);default:'enabled'"`
}

// TableName returns table name
func (SKU) TableName() string {
	return "shop_sku"
}
