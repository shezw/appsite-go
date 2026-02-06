// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package error

import (
	"errors"
	"testing"
)

func TestErrorCode(t *testing.T) {
	tests := []struct {
		code     ErrorCode
		expected int
		msg      string
	}{
		{Success, 200, "Success"},
		{InvalidParams, 400, "Invalid Parameters"},
		{ServerError, 500, "Internal Server Error"},
		{ErrorCode(999), 999, "Unknown Error"},
	}

	for _, tt := range tests {
		if tt.code.Int() != tt.expected {
			t.Errorf("ErrorCode.Int() = %d, want %d", tt.code.Int(), tt.expected)
		}
		if tt.code.String() != tt.msg {
			t.Errorf("ErrorCode.String() = %s, want %s", tt.code.String(), tt.msg)
		}
	}
}

func TestAppError(t *testing.T) {
	// Test New
	e1 := New(InvalidParams)
	if e1.Code != InvalidParams {
		t.Errorf("New() code = %v, want %v", e1.Code, InvalidParams)
	}
	if e1.Message != "Invalid Parameters" {
		t.Errorf("New() message = %v, want %v", e1.Message, "Invalid Parameters")
	}

	// Test NewWithMessage
	customMsg := "Custom invalid params"
	e2 := NewWithMessage(InvalidParams, customMsg)
	if e2.Message != customMsg {
		t.Errorf("NewWithMessage() message = %v, want %v", e2.Message, customMsg)
	}

	// Test Wrap
	baseErr := errors.New("db error")
	e3 := Wrap(ServerError, baseErr, "Database failed")
	if e3.Code != ServerError {
		t.Errorf("Wrap() code = %v, want %v", e3.Code, ServerError)
	}
	if e3.Err != baseErr {
		t.Errorf("Wrap() err = %v, want %v", e3.Err, baseErr)
	}
	if e3.Unwrap() != baseErr {
		t.Errorf("Unwrap() = %v, want %v", e3.Unwrap(), baseErr)
	}

	// Test WithDetails
	details := map[string]string{"field": "email"}
	e4 := New(InvalidParams).WithDetails(details)
	if e4.Details == nil {
		t.Error("WithDetails() failed to set details")
	}

	// Test Error() string
	expectedStr := "[500] Database failed: db error"
	if e3.Error() != expectedStr {
		t.Errorf("Error() = %s, want %s", e3.Error(), expectedStr)
	}

	e5 := New(NotFound)
	expectedStr2 := "[404] Not Found"
	if e5.Error() != expectedStr2 {
		t.Errorf("Error() = %s, want %s", e5.Error(), expectedStr2)
	}

	// Test Wrap with empty message
	e6 := Wrap(Unauthorized, nil, "")
	if e6.Message != Unauthorized.String() {
		t.Errorf("Wrap() empty msg = %v, want %v", e6.Message, Unauthorized.String())
	}
}
