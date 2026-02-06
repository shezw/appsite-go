// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package entity

import "appsite-go/internal/core/model"

// Deal represents a financial record (money or points)
type Deal struct {
	model.Base
	model.Tenant

	UserID      string `json:"user_id" gorm:"type:varchar(36);index;not null"`
	
	Type        string `json:"type" gorm:"type:varchar(32);index"` // payment, refund, reward, usage
	Asset       string `json:"asset" gorm:"type:varchar(16);default:'money'"` // money, point
	
	Amount      int64  `json:"amount" gorm:"comment:Positive for income, Negative for outcome"`
	Balance     int64  `json:"balance" gorm:"comment:Snapshot of balance after deal"`
	
	RelatedID   string `json:"related_id" gorm:"type:varchar(36);index"` // OrderID etc.
	Description string `json:"description" gorm:"type:varchar(255)"`
}

func (Deal) TableName() string {
	return "fin_deal"
}

// Balance holds the current state of a user's wallet
type Balance struct {
	model.Base
	model.Tenant // If balances are scoped by tenant

	UserID  string `json:"user_id" gorm:"type:varchar(36);uniqueIndex:idx_user_asset;not null"`
	Asset   string `json:"asset" gorm:"type:varchar(16);uniqueIndex:idx_user_asset;default:'money'"`
	
	Total   int64  `json:"total" gorm:"default:0"` // Current balance
	Frozen  int64  `json:"frozen" gorm:"default:0"` // Frozen amount
}

func (Balance) TableName() string {
	return "fin_balance"
}
