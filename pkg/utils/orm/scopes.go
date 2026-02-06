// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package orm

import "gorm.io/gorm"

// Paginate executes pagination by setting limit and offset.
// page is 1-based index (e.g., 1 is the first page).
// pageSize is the number of items per page.
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// Active returns a scope that filters by the 'state' column being 1 (active).
// Assumes the table has a 'state' column where 1 represents active.
func Active() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("state = ?", 1)
	}
}
