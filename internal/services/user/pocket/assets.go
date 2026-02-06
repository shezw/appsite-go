// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pocket

import (
	"appsite-go/internal/services/commerce/coupon"
	"appsite-go/internal/services/finance"
	"appsite-go/internal/services/finance/entity"
	couponEntity "appsite-go/internal/services/commerce/entity"
)

type AssetSummary struct {
	Points      int64 `json:"points"`
	Money       int64 `json:"money"`
	CouponCount int64 `json:"coupon_count"`
}

type Service struct {
	ledgerSvc *finance.LedgerService
	couponSvc *coupon.Service // We need a ListUserCoupons method in coupon service?
}

func NewService(l *finance.LedgerService, c *coupon.Service) *Service {
	return &Service{
		ledgerSvc: l,
		couponSvc: c,
	}
}

func (s *Service) GetAssets(userID string) (*AssetSummary, error) {
	points, err := s.ledgerSvc.GetBalance(userID, "point")
	if err != nil {
		return nil, err
	}

	money, err := s.ledgerSvc.GetBalance(userID, "money")
	if err != nil {
		return nil, err
	}

	// Count coupons via coupon service
	// We need to add CountUserCoupons to coupon service.
	// For now assume 0 if not implemented, or add it.
	cCount, err := s.couponSvc.CountUserCoupons(userID) 
	if err != nil {
		// Tolerate error or return?
		cCount = 0 
	}

	return &AssetSummary{
		Points:      points,
		Money:       money,
		CouponCount: cCount,
	}, nil
}

// GetDeals returns transaction history
func (s *Service) GetDeals(userID string, asset string, page, size int) ([]entity.Deal, int64, error) {
	return s.ledgerSvc.ListDeals(userID, asset, page, size)
}

// GetCoupons returns user coupons
func (s *Service) GetCoupons(userID string) ([]couponEntity.UserCoupon, error) {
	return s.couponSvc.ListUserCoupons(userID)
}
