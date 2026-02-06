// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stock

import (
	"errors"

	"gorm.io/gorm"

	"appsite-go/internal/services/commerce/entity"
)

var (
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidQuantity   = errors.New("invalid quantity")
)

// InventoryService handles stock operations
type InventoryService struct {
	db *gorm.DB
}

// NewInventoryService initializes the service
func NewInventoryService(db *gorm.DB) *InventoryService {
	if db != nil {
		_ = db.AutoMigrate(&entity.StockLog{})
	}
	return &InventoryService{db: db}
}

// Deduct reduces stock for an SKU. Safe concurrent update.
func (s *InventoryService) Deduct(skuID string, qty int, orderID string) error {
	if qty <= 0 {
		return ErrInvalidQuantity
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Check and Update atomically
		// GORM Update with expression
		// UPDATE shop_sku SET stock = stock - ?, updated_at = ? WHERE id = ? AND stock >= ?
		
		res := tx.Model(&entity.SKU{}).
			Where("id = ? AND stock >= ?", skuID, qty).
			Update("stock", gorm.Expr("stock - ?", qty))

		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			// Either ID not found or partial check failed (stock < qty)
			// We can check existence separately if we want distinct errors, 
			// but for high throughput, assume insufficient stock.
			return ErrInsufficientStock
		}

		// 2. Fetch updated SKU to log current state (optional, or just log delta)
		var sku entity.SKU
		if err := tx.Select("stock").First(&sku, "id = ?", skuID).Error; err != nil {
			return err
		}

		// 3. Record Log
		log := &entity.StockLog{
			SKUID:    skuID,
			OrderID:  orderID,
			Quantity: -qty, // Negative for deduct
			Type:     "order",
			After:    sku.Stock,
			Before:   sku.Stock + qty,
		}
		
		if err := tx.Create(log).Error; err != nil {
			return err
		}

		return nil
	})
}

// Restore increases stock (e.g., cancelled order)
func (s *InventoryService) Restore(skuID string, qty int, orderID string) error {
	if qty <= 0 {
		return ErrInvalidQuantity
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&entity.SKU{}).
			Where("id = ?", skuID).
			Update("stock", gorm.Expr("stock + ?", qty))
		
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("sku not found")
		}

		var sku entity.SKU
		if err := tx.Select("stock").First(&sku, "id = ?", skuID).Error; err != nil {
			return err
		}

		log := &entity.StockLog{
			SKUID:    skuID,
			OrderID:  orderID,
			Quantity: qty,
			Type:     "cancel",
			After:    sku.Stock,
			Before:   sku.Stock - qty,
		}

		if err := tx.Create(log).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetStock returns current stock
func (s *InventoryService) GetStock(skuID string) (int, error) {
	var stock int
	err := s.db.Model(&entity.SKU{}).Select("stock").Where("id = ?", skuID).Scan(&stock).Error
	return stock, err
}
