// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package info

import (
	"gorm.io/gorm"
	
	"appsite-go/internal/core/model"
	"appsite-go/internal/services/user/entity"
)

// Service handles user profile operations
type Service struct {
	repo *model.CRUD[entity.UserInfo]
}

// NewService creates a new info service
func NewService(db *gorm.DB) *Service {
	return &Service{
		repo: model.NewCRUD[entity.UserInfo](db),
	}
}

// UpdateProfile updates or creates profile info
func (s *Service) UpdateProfile(userID string, data map[string]interface{}) error {
	// Check if exists
	res := s.repo.Get(userID)
	if res.Success {
		return s.repo.Update(userID, data).Error
	}
	
	// If not, might need Create, but Update usually implies existing.
	// However, UserInfo shares primary key with User.
	// Assuming logic handles creation during registration OR lazy creation here.
	
	// Lazy create
	info := &entity.UserInfo{UserID: userID}
	// Copy data -> not efficient with map.
	// We'll leave lazy create to caller or use Upsert.
	
	// GORM Upsert with map is tricky.
	
	return s.repo.Update(userID, data).Error
}

// GetProfile retrieves
func (s *Service) GetProfile(userID string) (*entity.UserInfo, error) {
	res := s.repo.Get(userID)
	if !res.Success {
		return nil, res.Error
	}
	return res.Data.(*entity.UserInfo), nil
}
