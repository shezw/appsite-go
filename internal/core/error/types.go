// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package error

import (
	"fmt"
)

// AppError represents a structured error in the application
type AppError struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
	Err     error       `json:"-"` // underlying error, not exposed in JSON
}

// Error implements the standard error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError with a code and default message
func New(code ErrorCode) *AppError {
	return &AppError{
		Code:    code,
		Message: code.String(),
	}
}

// NewWithMessage creates a new AppError with a custom message
func NewWithMessage(code ErrorCode, msg string) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
	}
}

// Wrap creates a new AppError wrapping an existing error
// It uses the code's default message if msg is empty
func Wrap(code ErrorCode, err error, msg string) *AppError {
	if msg == "" {
		msg = code.String()
	}
	return &AppError{
		Code:    code,
		Message: msg,
		Err:     err,
	}
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details interface{}) *AppError {
	e.Details = details
	return e
}
