// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// AcquireLock attempts to acquire a distributed lock for the given key.
// It returns true if the lock was successfully acquired, false otherwise.
// value should be unique identifier for the owner (e.g. uuid/hostname+pid), used for safe unlocking.
func AcquireLock(ctx context.Context, rdb *redis.Client, key string, value string, expiration time.Duration) (bool, error) {
	return rdb.SetNX(ctx, key, value, expiration).Result()
}

// ReleaseLock releases the lock only if the value matches (CAS).
// This prevents removing a lock held by another client if the original lock expired.
func ReleaseLock(ctx context.Context, rdb *redis.Client, key string, value string) (bool, error) {
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`
	res, err := rdb.Eval(ctx, script, []string{key}, value).Result()
	if err != nil {
		return false, err
	}
	
	// res is int64 (0 or 1)
	if val, ok := res.(int64); ok {
		return val == 1, nil
	}
	return false, nil
}
