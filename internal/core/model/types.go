// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

// ListParams defines parameters for listing data
type ListParams struct {
	Page     int
	PageSize int
	Sort     string
	Filters  map[string]interface{}
	IsPublic bool
}

// Result mimics the ASResult structure for standardized returns
type Result struct {
	Success bool
	Data    interface{}
	Message string
	Error   error
}

// Operator defines the interface for model operations
type Operator[T any] interface {
	Add(entity *T) *Result
	Update(id string, updates map[string]interface{}) *Result
	Remove(id string) *Result
	Get(id string) *Result
	List(params *ListParams) *Result
}
