// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package finance

import (
	"errors"

	"appsite-go/internal/services/finance/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
)

type LedgerService struct {
	db *gorm.DB
}

func NewLedgerService(db *gorm.DB) *LedgerService {
	return &LedgerService{db: db}
}

func (s *LedgerService) Migrate() error {
	return s.db.AutoMigrate(&entity.Deal{}, &entity.Balance{})
}

// EnsureBalance creates a balance record if not exists
func (s *LedgerService) ensureBalance(tx *gorm.DB, userID string, asset string) error {
	var count int64
	tx.Model(&entity.Balance{}).Where("user_id = ? AND asset = ?", userID, asset).Count(&count)
	if count == 0 {
		return tx.Create(&entity.Balance{
			UserID: userID,
			Asset:  asset,
			Total:  0,
		}).Error
	}
	return nil
}

// RecordTransaction updates balance and logs a deal
func (s *LedgerService) RecordTransaction(userID string, asset string, amount int64, dealType string, relatedID string, desc string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Ensure balance record
		if err := s.ensureBalance(tx, userID, asset); err != nil {
			return err
		}

		// 2. Lock and Update Balance
		var bal entity.Balance
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND asset = ?", userID, asset).
			First(&bal).Error; err != nil {
			return err
		}

		newTotal := bal.Total + amount
		if newTotal < 0 {
			return ErrInsufficientFunds
		}

		if err := tx.Model(&bal).Update("total", newTotal).Error; err != nil {
			return err
		}

		// 3. Create Deal record
		deal := &entity.Deal{
			UserID:      userID,
			Asset:       asset,
			Type:        dealType,
			Amount:      amount,
			Balance:     newTotal,
			RelatedID:   relatedID,
			Description: desc,
		}
		return tx.Create(deal).Error
	})
}

// GetBalance returns current balance
func (s *LedgerService) GetBalance(userID string, asset string) (int64, error) {
	var bal entity.Balance
	err := s.db.Where("user_id = ? AND asset = ?", userID, asset).First(&bal).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	return bal.Total, err
}

// ListDeals returns history
func (s *LedgerService) ListDeals(userID string, asset string, page, size int) ([]entity.Deal, int64, error) {
	var deals []entity.Deal
	var total int64
	
	db := s.db.Model(&entity.Deal{}).Where("user_id = ?", userID)
	if asset != "" {
		db = db.Where("asset = ?", asset)
	}

	db.Count(&total)
	
	err := db.Order("created_at desc").Offset((page - 1) * size).Limit(size).Find(&deals).Error
	return deals, total, err
}
