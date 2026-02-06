// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package account_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	
	"appsite-go/internal/core/setting"
	"appsite-go/internal/services/access/token"
	"appsite-go/internal/services/access/verify"
	"appsite-go/internal/services/user/account"
)


func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func setupAuthComponents(t *testing.T, db *gorm.DB) (*account.AuthService, *verify.OTPService) {
	// Redis setup
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(s.Close)
	rdb := redis.NewClient(&redis.Options{Addr: s.Addr()})
	otpSvc := verify.NewOTPService(rdb)

	cfg := setting.AppConfig{
		JwtSecret: "test_secret",
		JwtExpire: time.Hour,
	}
	tokenSvc := token.NewService(cfg)
	return account.NewAuthService(db, tokenSvc, otpSvc), otpSvc
}

func TestRegisterAndLogin(t *testing.T) {
	db := setupDB(t)
	svc, _ := setupAuthComponents(t, db)

	input := account.RegisterInput{
		Username: "alice",
		Password: "password123",
		Email:    "alice@example.com",
		Nickname: "Alice Wonderland",
	}

	// 1. Register
	user, err := svc.Register(input)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	if user.ID == "" {
		t.Error("User ID not generated")
	}
	// Password should be hashed
	if user.Password == input.Password {
		t.Error("Password storage cleartext detected")
	}

	// 2. Duplicate Check
	_, err = svc.Register(input)
	if err != account.ErrUserExists {
		t.Errorf("Duplicate register should fail with ErrUserExists, got %v", err)
	}

	// 3. Login
	tokenStr, loginUser, err := svc.Login("alice", "password123")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if tokenStr == "" {
		t.Error("Token empty")
	}
	if loginUser.ID != user.ID {
		t.Error("Login returned wrong user")
	}

	// 4. Login Fail
	_, _, err = svc.Login("alice", "wrongpass")
	if err != account.ErrInvalidPwd {
		t.Error("Expected ErrInvalidPwd")
	}

	_, _, err = svc.Login("bob", "password123")
	if err != account.ErrUserNotFound {
		t.Error("Expected ErrUserNotFound")
	}

	// 5. Disabled User
	user.Status = "disabled"
	db.Save(user)
	
	_, _, err = svc.Login("alice", "password123")
	if err != account.ErrUserDisabled {
		t.Error("Expected ErrUserDisabled")
	}
}

func TestOTPLogin(t *testing.T) {
	db := setupDB(t)
	svc, otp := setupAuthComponents(t, db)
	ctx := context.Background()

	// Register a user
	svc.Register(account.RegisterInput{
		Mobile: "13800000000",
		Password: "password",
	})
	
	// Generate OTP
	code, _ := otp.Generate(ctx, "13800000000", 6, time.Minute)
	
	// Valid Access
	_, _, err := svc.LoginByOTP(ctx, "13800000000", code)
	if err != nil {
		t.Errorf("OTP Login failed: %v", err)
	}
	
	// Invalid Code
	_, _, err = svc.LoginByOTP(ctx, "13800000000", "000000")
	if err != account.ErrInvalidOTP {
		t.Errorf("Expected invalid OTP, got %v", err)
	}
}

func TestRegisterMobile(t *testing.T) {
	db := setupDB(t)
	svc, _ := setupAuthComponents(t, db)
	
	u, err := svc.RegisterByMobile("13912345678", "pass")
	if err != nil {
		t.Fatal(err)
	}
	if u.Username != "m_13912345678" {
		t.Error("Username generation failed")
	}
}
