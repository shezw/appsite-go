// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package account

import "golang.org/x/crypto/bcrypt"

// PasswordService handles hashing and verification
type PasswordService struct {
    cost int
}

// NewPasswordService creates a new password service
func NewPasswordService() *PasswordService {
    return &PasswordService{
        cost: bcrypt.DefaultCost,
    }
}

// Hash creates a bcrypt hash of the password
func (s *PasswordService) Hash(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
    return string(bytes), err
}

// Compare checks if the password matches the hash
func (s *PasswordService) Compare(hash, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
