// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package entity

import "appsite-go/internal/core/model"

// Coupon represents a discount rule template
type Coupon struct {
	model.Base
	model.Tenant

	Title       string `json:"title" gorm:"type:varchar(64);not null"`
	Description string `json:"description" gorm:"type:varchar(255)"`
	
	Type        string `json:"type" gorm:"type:varchar(16);comment:cash, discount"` // cash (fixed), discount (percentage)
	Value       int64  `json:"value" gorm:"comment:Amount in cents or Percentage (e.g. 80 for 80%)"`
	
	MinSpend    int64  `json:"min_spend" gorm:"default:0;comment:Minimum order amount in cents"`
	
	StartTime   int64  `json:"start_time"`
	EndTime     int64  `json:"end_time"`
	
	TotalCount  int    `json:"total_count" gorm:"default:-1;comment:-1 unlimited"`
	TakenCount  int    `json:"taken_count" gorm:"default:0"`
	
	Status      string `json:"status" gorm:"type:varchar(16);default:'enabled'"`
}

// TableName returns table name
func (Coupon) TableName() string {
	return "shop_coupon"
}

// UserCoupon represents a coupon instance held by a user
type UserCoupon struct {
	model.Base
	
	UserID    string `json:"user_id" gorm:"type:varchar(36);index;not null"`
	CouponID  string `json:"coupon_id" gorm:"type:varchar(36);index;not null"`
	
	Status    string `json:"status" gorm:"type:varchar(16);default:'unused';index"` // unused, used, expired
	UsedAt    int64  `json:"used_at"`
	OrderID   string `json:"order_id" gorm:"type:varchar(36);index"` // If used
}

// TableName returns table name
func (UserCoupon) TableName() string {
	return "shop_user_coupon"
}
