// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package account

import (
	"context"
	"errors"

	"gorm.io/gorm"
	
	"appsite-go/internal/core/model"
	"appsite-go/internal/services/access/token"
	"appsite-go/internal/services/access/verify"
	"appsite-go/internal/services/user/entity"
)

var (
	ErrUserExists    = errors.New("user already exists")
	ErrUserNotFound  = errors.New("user not found")
	ErrInvalidPwd    = errors.New("invalid password")
	ErrUserDisabled  = errors.New("user is disabled")
	ErrInvalidOTP    = errors.New("invalid otp code")
)

// AuthService handles authentication
type AuthService struct {
	db       *gorm.DB
	repo     *model.CRUD[entity.User]
	pwd      *PasswordService
	tokenSvc *token.Service
	otpSvc   *verify.OTPService
}

// NewAuthService creates a new auth service
func NewAuthService(db *gorm.DB, tokenSvc *token.Service, otpSvc *verify.OTPService) *AuthService {
	// Auto migrate
	if db != nil {
		_ = db.AutoMigrate(&entity.User{}, &entity.UserInfo{}, &entity.UserGroup{}, &entity.UserPreference{})
	}
	return &AuthService{
		db:       db,
		repo:     model.NewCRUD[entity.User](db),
		pwd:      NewPasswordService(),
		tokenSvc: tokenSvc,
		otpSvc:   otpSvc,
	}
}

// RegisterInput defines parameters for registration
type RegisterInput struct {
	Username string
	Password string
	Email    string
	Mobile   string
	Nickname string
}

// Register creates a new user user
func (s *AuthService) Register(input RegisterInput) (*entity.User, error) {
	// 1. Check existence (Email or Username)
	var count int64
	s.db.Model(&entity.User{}).
		Where("username = ? OR email = ?", input.Username, input.Email).
		Count(&count)
	
	if count > 0 {
		return nil, ErrUserExists
	}

	// 2. Hash Password
	hashedPwd, err := s.pwd.Hash(input.Password)
	if err != nil {
		return nil, err
	}

	// 3. Create
	user := &entity.User{
		Username: input.Username,
		Email:    input.Email,
		Mobile:   input.Mobile,
		Password: hashedPwd,
		Nickname: input.Nickname,
		Status:   "enabled",
	}

	if res := s.repo.Add(user); !res.Success {
		// Unwrap error if possible
		if res.Error != nil {
			return nil, res.Error
		}
		return nil, errors.New(res.Message)
	}

	return user, nil
}

// Login verifies credentials and returns a token
func (s *AuthService) Login(identifier, password string) (string, *entity.User, error) {
	user := &entity.User{}
	
	// Find (support username/email/mobile login)
	err := s.db.Where("username = ? OR email = ? OR mobile = ?", identifier, identifier, identifier).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, ErrUserNotFound
		}
		return "", nil, err
	}

	// Verify Password
	if !s.pwd.Compare(user.Password, password) {
		return "", nil, ErrInvalidPwd
	}

	// Check Status
	if user.Status != "enabled" {
		return "", nil, ErrUserDisabled
	}

	// Generate Token
	tokenStr, err := s.tokenSvc.GenerateToken(user.ID, user.GroupID)
	if err != nil {
		return "", nil, err
	}

	return tokenStr, user, nil
}

// LoginByOTP logs in using mobile/email + OTP
func (s *AuthService) LoginByOTP(ctx context.Context, target, code string) (string, *entity.User, error) {
	// 1. Verify OTP
	if !s.otpSvc.Check(ctx, target, code) {
		return "", nil, ErrInvalidOTP
	}

	// 2. Find User (Mobile or Email)
	user := &entity.User{}
	err := s.db.Where("mobile = ? OR email = ?", target, target).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Optional: Auto-register logic could be placed here? 
			// For now, strict login
			return "", nil, ErrUserNotFound
		}
		return "", nil, err
	}

	// 3. Status Check
	if user.Status != "enabled" {
		return "", nil, ErrUserDisabled
	}

	// 4. Token
	tokenStr, err := s.tokenSvc.GenerateToken(user.ID, user.GroupID)
	if err != nil {
		return "", nil, err
	}

	return tokenStr, user, nil
}

// RegisterByMobile registers a new user with mobile
func (s *AuthService) RegisterByMobile(mobile, password string) (*entity.User, error) {
	return s.Register(RegisterInput{
		Mobile:   mobile,
		Password: password,
		Username: "m_" + mobile, // Auto-gen username
		Nickname: "User_" + mobile[len(mobile)-4:], // Suffix
	})
}
