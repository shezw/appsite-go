// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package entity

import "appsite-go/internal/core/model"

type Order struct {
	model.Base
	model.Tenant

	OrderNo     string `json:"order_no" gorm:"type:varchar(32);uniqueIndex;not null"`
	UserID      string `json:"user_id" gorm:"type:varchar(36);index;not null"`
	
	TotalAmount int64  `json:"total_amount" gorm:"comment:Sum of item prices"`
	PayAmount   int64  `json:"pay_amount" gorm:"comment:Final amount to pay"`
	Discount    int64  `json:"discount" gorm:"comment:Discount applied"`
	
	Status      string `json:"status" gorm:"type:varchar(16);index;default:'pending'"` 
	// pending, paid, shipping, done, closed, refunded

	// Payment Info
	PayMethod   string `json:"pay_method" gorm:"type:varchar(16)"` // wechat, alipay, stripe
	TransactionID string `json:"transaction_id" gorm:"type:varchar(64);index"` // 3rd party ID

	// Snapshot Data (Using simple string/json for now, could be JSON types in Postgres)
	AddressSnapshot string `json:"address_snapshot" gorm:"type:text"`
	Note            string `json:"note" gorm:"type:varchar(255)"`
}

func (Order) TableName() string {
	return "shop_order"
}

type OrderItem struct {
	model.Base
	
	OrderID     string `json:"order_id" gorm:"type:varchar(36);index;not null"`
	
	ProductID   string `json:"product_id" gorm:"type:varchar(36);index"`
	SkuID       string `json:"sku_id" gorm:"type:varchar(36);index"`
	
	Title       string `json:"title" gorm:"type:varchar(128)"`
	SkuSpec     string `json:"sku_spec" gorm:"type:varchar(255)"` // e.g. "Color: Red, Size: L"
	Thumb       string `json:"thumb" gorm:"type:varchar(255)"`
	
	Price       int64  `json:"price" gorm:"comment:Unit price at moment of purchase"`
	Quantity    int    `json:"quantity"`
	Amount      int64  `json:"amount" gorm:"comment:Price * Quantity"`
}

func (OrderItem) TableName() string {
	return "shop_order_item"
}
