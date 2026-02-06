// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package writeoff

import (
	"errors"

	"appsite-go/internal/services/commerce/coupon"
	"appsite-go/internal/services/commerce/entity"
)

var (
	ErrInvalidCode = errors.New("invalid code or expired")
)

type Service struct {
	couponSvc *coupon.Service
	// orderSvc *order.Service // In future could verify tickets/pickup orders
}

func NewService(cSvc *coupon.Service) *Service {
	return &Service{
		couponSvc: cSvc,
	}
}

// VerifyCouponCode checks if a coupon code (ID) is valid for write-off.
// Unlike rule.Verify which checks for 'Apply' (min spend etc), this checks for 'Redemption' (existence, status).
func (s *Service) VerifyCouponCode(code string) (*entity.UserCoupon, *entity.Coupon, error) {
	uc, err := s.couponSvc.GetUserCoupon(code)
	if err != nil {
		return nil, nil, err
	}
	
	if uc.Status != "unused" {
		return nil, nil, ErrInvalidCode
	}

	c, err := s.couponSvc.GetCoupon(uc.CouponID)
	if err != nil {
		return nil, nil, err
	}
	
	return uc, c, nil
}

// WriteOffCoupon consumes the coupon.
func (s *Service) WriteOffCoupon(code string, staffID string) error {
	// Record who wrote it off?
	// UserCoupon schema doesn't have 'WrittenOffBy'.
	// Maybe just mark used.
	return s.couponSvc.Use(code, "writeoff:"+staffID)
}
