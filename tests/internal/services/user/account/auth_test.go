// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package account_test

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	
	"appsite-go/internal/core/setting"
	"appsite-go/internal/services/access/token"
	"appsite-go/internal/services/user/account"
)

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func setupAuthService(t *testing.T, db *gorm.DB) *account.AuthService {
	cfg := setting.AppConfig{
		JwtSecret: "test_secret",
		JwtExpire: time.Hour,
	}
	tokenSvc := token.NewService(cfg)
	return account.NewAuthService(db, tokenSvc)
}

func TestRegisterAndLogin(t *testing.T) {
	db := setupDB(t)
	svc := setupAuthService(t, db)

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
