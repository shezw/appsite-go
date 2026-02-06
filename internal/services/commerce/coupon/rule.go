// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package coupon

import (
	"errors"
	"time"

	"appsite-go/internal/services/commerce/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrCouponNotFound = errors.New("coupon not found")
	ErrCouponExpired  = errors.New("coupon expired or not started")
	ErrCouponEmpty    = errors.New("coupon out of stock")
	ErrAlreadyTaken   = errors.New("limit reached for this user")
	ErrMinSpend       = errors.New("minimum spend not met")
	ErrCouponUsed     = errors.New("coupon already used or invalid")
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Migrate() error {
	return s.db.AutoMigrate(&entity.Coupon{}, &entity.UserCoupon{})
}

// CreateCoupon creates a new coupon rule
func (s *Service) CreateCoupon(c *entity.Coupon) error {
	return s.db.Create(c).Error
}

// Issue grants a coupon to a user
func (s *Service) Issue(getUserID func() string, couponID string) (*entity.UserCoupon, error) {
	userId := getUserID()
	now := time.Now().Unix()

	var uc *entity.UserCoupon

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var c entity.Coupon
		// Lock the coupon row for atomic counter update
		if err := tx.Clauses(clauseLocking).First(&c, "id = ?", couponID).Error; err != nil {
			return err
		}

		if c.Status != "enabled" {
			return ErrCouponNotFound
		}

		if now < c.StartTime || now > c.EndTime {
			return ErrCouponExpired
		}

		if c.TotalCount > -1 && c.TakenCount >= c.TotalCount {
			return ErrCouponEmpty
		}

		// Check if user already has this coupon (simple restriction: 1 per user for now, or you could add a MaxPerUser field to Coupon)
		// For simplicity in this iteration, we assume 1 per user per coupon type to avoid spam.
		var count int64
		tx.Model(&entity.UserCoupon{}).Where("user_id = ? AND coupon_id = ?", userId, couponID).Count(&count)
		if count > 0 {
			return ErrAlreadyTaken
		}

		// Update taken count
		if err := tx.Model(&c).Update("taken_count", gorm.Expr("taken_count + 1")).Error; err != nil {
			return err
		}

		// Create UserCoupon
		uc = &entity.UserCoupon{
			UserID:   userId,
			CouponID: couponID,
			Status:   "unused",
		}
		return tx.Create(uc).Error
	})

	return uc, err
}

// Verify checks if a user coupon is valid for a given order amount
func (s *Service) Verify(userCouponID string, orderAmount int64) (*entity.Coupon, error) {
	var uc entity.UserCoupon
	if err := s.db.First(&uc, "id = ?", userCouponID).Error; err != nil {
		return nil, ErrCouponNotFound
	}

	if uc.Status != "unused" {
		return nil, ErrCouponUsed
	}

	var c entity.Coupon
	if err := s.db.First(&c, "id = ?", uc.CouponID).Error; err != nil {
		return nil, ErrCouponNotFound
	}

	now := time.Now().Unix()
	if now > c.EndTime {
		return nil, ErrCouponExpired
	}

	if orderAmount < c.MinSpend {
		return nil, ErrMinSpend
	}

	return &c, nil
}

// GetUserCoupon retrieves raw user coupon data
func (s *Service) GetUserCoupon(id string) (*entity.UserCoupon, error) {
	var uc entity.UserCoupon
	if err := s.db.First(&uc, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &uc, nil
}

// GetCoupon retrieves raw coupon definition
func (s *Service) GetCoupon(id string) (*entity.Coupon, error) {
	var c entity.Coupon
	if err := s.db.First(&c, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

// CountUserCoupons returns number of unused coupons
func (s *Service) CountUserCoupons(userID string) (int64, error) {
	var count int64
	err := s.db.Model(&entity.UserCoupon{}).
		Where("user_id = ? AND status = ?", userID, "unused").
		Count(&count).Error
	return count, err
}

// ListUserCoupons returns list of unused coupons
func (s *Service) ListUserCoupons(userID string) ([]entity.UserCoupon, error) {
	var list []entity.UserCoupon
	err := s.db.Where("user_id = ? AND status = ?", userID, "unused").Find(&list).Error
	return list, err
}

// Use marks a coupon as used
func (s *Service) Use(userCouponID string, orderID string) error {
	res := s.db.Model(&entity.UserCoupon{}).
		Where("id = ? AND status = ?", userCouponID, "unused").
		Updates(map[string]interface{}{
			"status":   "used",
			"used_at":  time.Now().Unix(),
			"order_id": orderID,
		})
	
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrCouponUsed
	}
	return nil
}

// Helper for locking used in Issue
// We need to define clauseLocking safely. 
// Standard GORM clause.Locking{Strength: "UPDATE"} works for Postgres/MySQL.
// For SQLite, it's ignored but harmless, or we can look it up.
// Let's use a raw query logic or safe generic provided by GORM.

var clauseLocking = clause.Locking{Strength: "UPDATE"}
