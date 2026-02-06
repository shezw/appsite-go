// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package saas

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"appsite-go/internal/core/model"
	"appsite-go/internal/services/world/entity"
)

// TenantService handles SaaS tenant operations
type TenantService struct {
	db   *gorm.DB
	repo *model.CRUD[entity.Tenant]
}

// NewTenantService initializes the service
func NewTenantService(db *gorm.DB) *TenantService {
	if db != nil {
		_ = db.AutoMigrate(&entity.Tenant{})
	}
	return &TenantService{
		db:   db,
		repo: model.NewCRUD[entity.Tenant](db),
	}
}

// Create registers a new tenant
func (s *TenantService) Create(tenant *entity.Tenant) error {
	// Validate uniqueness of Domain/Code logic if needed beyond DB constraints
	if tenant.Code == "" {
		return errors.New("tenant code is required")
	}
	res := s.repo.Add(tenant)
	return res.Error
}

// Update modifies tenant info
func (s *TenantService) Update(id string, updates map[string]interface{}) error {
	res := s.repo.Update(id, updates)
	return res.Error
}

// Get retrieves a tenant by ID
func (s *TenantService) Get(id string) (*entity.Tenant, error) {
	res := s.repo.Get(id)
	if !res.Success {
		return nil, res.Error
	}
	return res.Data.(*entity.Tenant), nil
}

// GetByDomain retrieves a tenant by domain or code
func (s *TenantService) GetByDomain(domain string) (*entity.Tenant, error) {
	var t entity.Tenant
	err := s.db.Where("domain = ? OR code = ?", domain, domain).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// List returns tenants
func (s *TenantService) List(page, size int, filters map[string]interface{}) ([]entity.Tenant, int64, error) {
	res := s.repo.List(&model.ListParams{
		Page:     page,
		PageSize: size,
		Filters:  filters,
		Sort:     "created_at desc",
	})
	
	if !res.Success {
		return nil, 0, res.Error
	}
	
	data := res.Data.(map[string]interface{})
	return data["list"].([]entity.Tenant), data["total"].(int64), nil
}

// CheckActive verifies if a tenant is active and not expired
func (s *TenantService) CheckActive(id string) (bool, error) {
	t, err := s.Get(id)
	if err != nil {
		return false, err
	}
	
	if t.Status != "enabled" {
		return false, nil
	}
	
	if t.ExpireAt > 0 && t.ExpireAt < time.Now().Unix() {
		return false, nil
	}
	
	return true, nil
}
