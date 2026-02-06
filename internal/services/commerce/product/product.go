// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package product

import (
	"gorm.io/gorm"

	"appsite-go/internal/core/model"
	"appsite-go/internal/services/commerce/entity"
)

// Service handles product operations
type Service struct {
	db      *gorm.DB
	repo    *model.CRUD[entity.Product]
	skuRepo *model.CRUD[entity.SKU]
}

// NewService initializes the service
func NewService(db *gorm.DB) *Service {
	if db != nil {
		_ = db.AutoMigrate(&entity.Product{}, &entity.SKU{})
	}
	return &Service{
		db:      db,
		repo:    model.NewCRUD[entity.Product](db),
		skuRepo: model.NewCRUD[entity.SKU](db),
	}
}

// CreateProduct adds a new SPU
func (s *Service) CreateProduct(p *entity.Product) error {
	res := s.repo.Add(p)
	return res.Error
}

// UpdateProduct modifies SPU
func (s *Service) UpdateProduct(id string, updates map[string]interface{}) error {
	res := s.repo.Update(id, updates)
	return res.Error
}

// GetProduct retrieves SPU
func (s *Service) GetProduct(id string) (*entity.Product, error) {
	res := s.repo.Get(id)
	if !res.Success {
		return nil, res.Error
	}
	return res.Data.(*entity.Product), nil
}

// CreateSKU adds a new SKU
func (s *Service) CreateSKU(sku *entity.SKU) error {
	res := s.skuRepo.Add(sku)
	return res.Error
}

// UpdateSKU modifies SKU (Price, Stock, etc)
func (s *Service) UpdateSKU(id string, updates map[string]interface{}) error {
	res := s.skuRepo.Update(id, updates)
	return res.Error
}

// GetSKU retrieves SKU
func (s *Service) GetSKU(id string) (*entity.SKU, error) {
	res := s.skuRepo.Get(id)
	if !res.Success {
		return nil, res.Error
	}
	return res.Data.(*entity.SKU), nil
}

// ListSkus returns all SKUs for a product
func (s *Service) ListSkus(productID string) ([]entity.SKU, error) {
	var skus []entity.SKU
	err := s.db.Where("product_id = ?", productID).Find(&skus).Error
	return skus, err
}

// DeleteProduct removes product and its SKUs (Logical delete via CRUD)
func (s *Service) DeleteProduct(id string) error {
	// Use transaction to ensure both are deleted.
	// We use tx directly instead of repo to ensure we stay in the same transaction connection.
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Verify existence if needed, or just delete.
		// Note: This skips CRUD hooks if any.
		if err := tx.Delete(&entity.Product{}, "id = ?", id).Error; err != nil {
			return err
		}
		// Delete SKUs
		if err := tx.Where("product_id = ?", id).Delete(&entity.SKU{}).Error; err != nil {
			return err
		}
		return nil
	})
}
