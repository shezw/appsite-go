// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package verify_test

import (
	"context"
	"testing"
	"time"

	"appsite-go/internal/services/access/verify"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func setupRedis(t *testing.T) *redis.Client {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	t.Cleanup(s.Close)

	return redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
}

func TestOTP_Flow(t *testing.T) {
	rdb := setupRedis(t)
	svc := verify.NewOTPService(rdb)
	ctx := context.Background()

	target := "user@example.com"

	// 1. Generate
	code, err := svc.Generate(ctx, target, 6, time.Minute)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if len(code) != 6 {
		t.Errorf("Expected 6 chars, got %d", len(code))
	}

	// 2. Check Fail
	if svc.Check(ctx, target, "000000") {
		t.Error("Check should fail with wrong code")
	}

	// 3. Check Success
	if !svc.Check(ctx, target, code) {
		t.Error("Check should pass with correct code")
	}

	// 4. Replay Attack (Check again should fail)
	if svc.Check(ctx, target, code) {
		t.Error("OTP should be consumed after use")
	}
}

func TestOTP_Expiration(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	rdb := redis.NewClient(&redis.Options{Addr: s.Addr()})
	svc := verify.NewOTPService(rdb)
	ctx := context.Background()

	code, _ := svc.Generate(ctx, "expire@me", 4, time.Millisecond*100)
	
	s.FastForward(time.Second)
	
	if svc.Check(ctx, "expire@me", code) {
		t.Error("Expired OTP should verify false")
	}
}

func TestOTP_InvalidArgs(t *testing.T) {
	svc := verify.NewOTPService(nil)
	_, err := svc.Generate(context.Background(), "t", 0, time.Second)
	if err == nil {
		t.Error("Should fail on lengths <= 0")
	}
}
