// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package group

import (
	"gorm.io/gorm"
	
	"appsite-go/internal/core/model"
	"appsite-go/internal/services/user/entity"
)

type Service struct {
	repo *model.CRUD[entity.UserGroup]
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		repo: model.NewCRUD[entity.UserGroup](db),
	}
}

func (s *Service) Create(name, desc string, level int) (*entity.UserGroup, error) {
	g := &entity.UserGroup{
		GroupName:   name,
		Description: desc,
		Level:       level,
	}
	res := s.repo.Add(g)
	return g, res.Error
}

func (s *Service) List() ([]entity.UserGroup, error) {
	// Simple list
	res := s.repo.List(&model.ListParams{})
	if !res.Success {
		return nil, res.Error
	}
	data := res.Data.(map[string]interface{})
	return data["list"].([]entity.UserGroup), nil
}
