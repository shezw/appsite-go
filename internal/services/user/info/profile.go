// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package info

import (
	"errors"

	"gorm.io/gorm"

	"appsite-go/internal/core/model"
	"appsite-go/internal/services/user/entity"
)

// Service handles user profile operations
type Service struct {
	db   *gorm.DB
	repo *model.CRUD[entity.UserInfo]
}

// NewService creates a new info service
func NewService(db *gorm.DB) *Service {
	return &Service{
		db:   db,
		repo: model.NewCRUD[entity.UserInfo](db),
	}
}

// UpdateProfile updates or creates profile info
func (s *Service) UpdateProfile(userID string, data map[string]interface{}) error {
	var info entity.UserInfo
	err := s.db.First(&info, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			data["user_id"] = userID
			return s.db.Model(&entity.UserInfo{}).Create(data).Error
		}
		return err
	}
	return s.db.Model(&info).Updates(data).Error
}

// GetProfile retrieves
func (s *Service) GetProfile(userID string) (*entity.UserInfo, error) {
	var info entity.UserInfo
	if err := s.db.First(&info, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &info, nil
}
