// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package entity

import "appsite-go/internal/core/model"

// StockLog records inventory changes
type StockLog struct {
	model.Base
	
	SKUID    string `json:"sku_id" gorm:"column:sku_id;type:varchar(36);index;not null"`
	OrderID  string `json:"order_id" gorm:"column:order_id;type:varchar(36);index;comment:Order ID causing change"`
	
	Quantity int    `json:"quantity" gorm:"column:quantity;comment:Positive for add, Negative for deduct"`
	Type     string `json:"type" gorm:"column:type;type:varchar(16)"` // "order", "cancel", "admin", "return"
	
	Before   int    `json:"before" gorm:"column:before;comment:Snapshot before change"`
	After    int    `json:"after" gorm:"column:after;comment:Snapshot after change"`
}

// TableName returns table name
func (StockLog) TableName() string {
	return "shop_stock_log"
}
