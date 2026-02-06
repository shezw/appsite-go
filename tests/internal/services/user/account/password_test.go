// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package account_test

import (
	"testing"

	"appsite-go/internal/services/user/account"
)

func TestPassword(t *testing.T) {
	svc := account.NewPasswordService()
	pwd := "mySecret123"

	// Hash
	hash, err := svc.Hash(pwd)
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}
	if hash == pwd {
		t.Error("Hash should not match plaintext")
	}

	// Compare Success
	if !svc.Compare(hash, pwd) {
		t.Error("Compare failed for correct password")
	}

	// Compare Fail
	if svc.Compare(hash, "wrong") {
		t.Error("Compare passed for wrong password")
	}
}
