// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package verify

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/redis/go-redis/v9"
)

// OTPService handles One-Time Password generation and verification
type OTPService struct {
	rdb *redis.Client
}

// NewOTPService creates a new OTP service instance
func NewOTPService(rdb *redis.Client) *OTPService {
	return &OTPService{rdb: rdb}
}

// Generate creates a numeric OTP code and stores it in Redis
// target: unique identifier (email, phone)
// length: number of digits (e.g. 6)
// ttl: validity duration
func (s *OTPService) Generate(ctx context.Context, target string, length int, ttl time.Duration) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid length")
	}

	code := ""
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code += n.String()
	}

	key := s.key(target)
	if err := s.rdb.Set(ctx, key, code, ttl).Err(); err != nil {
		return "", err
	}

	return code, nil
}

// Check validates the OTP code. Returns true if valid and deletes the key.
func (s *OTPService) Check(ctx context.Context, target, code string) bool {
	key := s.key(target)
	val, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		return false // Key doesn't exist or redis error
	}

	if val == code {
		// Valid, consume it
		s.rdb.Del(ctx, key)
		return true
	}

	return false
}

func (s *OTPService) key(target string) string {
	return fmt.Sprintf("verify:otp:%s", target)
}
