// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package preference

import (
	"gorm.io/gorm"
	
	"appsite-go/internal/core/model"
	"appsite-go/internal/services/user/entity"
)

type Service struct {
	db   *gorm.DB
	repo *model.CRUD[entity.UserPreference]
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db:   db,
		repo: model.NewCRUD[entity.UserPreference](db),
	}
}

// Set uploads or updates a preference key
func (s *Service) Set(userID, saasID, key, content string) error {
	var pref entity.UserPreference
	err := s.db.Where("user_id = ? AND key_id = ?", userID, key).First(&pref).Error
	
	if err == nil {
		// Update
		return s.repo.Update(pref.ID, map[string]interface{}{"Content": content}).Error
	}
	
	if err == gorm.ErrRecordNotFound {
		// Create
		return s.repo.Add(&entity.UserPreference{
			UserID:  userID,
			SaasID:  saasID,
			KeyID:   key,
			Content: content,
		}).Error
	}
	
	return err
}

// Get retrieves content
func (s *Service) Get(userID, key string) (string, error) {
	var pref entity.UserPreference
	err := s.db.Where("user_id = ? AND key_id = ?", userID, key).First(&pref).Error
	if err != nil {
		return "", err
	}
	return pref.Content, nil
}
