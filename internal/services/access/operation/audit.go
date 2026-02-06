// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package operation

import (
	"appsite-go/internal/core/model"

	"gorm.io/gorm"
)

// AuditLog records user activities
type AuditLog struct {
	model.Base
	UserID    string `gorm:"index;size:32"`
	TenantID  string `gorm:"index;size:32"`
	Action    string `gorm:"size:64;index"` // e.g. "login", "create_post"
	Method    string `gorm:"size:10"`       // HTTP Method
	Path      string `gorm:"size:255"`      // Request Path
	IP        string `gorm:"size:45"`
	UserAgent string `gorm:"size:255"`
	Status    int    // HTTP Status
	Detail    string `gorm:"type:text"` // JSON payload or error message
}

// Service handles audit logging
type Service struct {
	db *gorm.DB
}

// NewService creates a new audit service
func NewService(db *gorm.DB) *Service {
	// Auto migrate
	if db != nil {
		_ = db.AutoMigrate(&AuditLog{})
	}
	return &Service{db: db}
}

// Record saves an audit log entry
func (s *Service) Record(log *AuditLog) error {
	return s.db.Create(log).Error
}

// FindByUser retrieves logs for a user
func (s *Service) FindByUser(userID string, limit int) ([]AuditLog, error) {
	var logs []AuditLog
	err := s.db.Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}
