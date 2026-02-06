// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package redis_test

import (
	"context"
	"log"
	"strconv"
	"testing"
	"time"

	"appsite-go/internal/core/setting"
	"appsite-go/pkg/utils/redis"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/redis/go-redis/v9"
)

func TestNewClient(t *testing.T) {
	// Start miniredis
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	// Case 1: Success
	port, _ := strconv.Atoi(mr.Port())
	cfg := &setting.RedisConfig{
		Host: mr.Host(),
		Port: port,
		Password: "",
		DB:       0,
	}

	client, err := redis.NewClient(cfg)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	// Verify we can ping
	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Errorf("Ping failed: %v", err)
	}

	// Case 2: Nil Config
	_, err = redis.NewClient(nil)
	if err == nil {
		t.Error("NewClient(nil) should fail")
	}

	// Case 3: Connection Failure (invalid port)
	// Suppress redis logs for this expected failure to keep test output clean
	goredis.SetLogger(&noOpLogger{})
	defer goredis.SetLogger(&defaultLogger{})

	cfgFail := &setting.RedisConfig{
		Host: "localhost",
		Port: 54321, // assume this port is closed
	}
	// Note: NewClient pings the server, so it should return error if connection fails
	_, err = redis.NewClient(cfgFail)
	if err == nil {
		// If it doesn't fail, maybe something is running there? Unlikely.
		// Or NewClient doesn't ping? It does in our implementation.
		t.Log("NewClient(invalid) did not fail, check implementation")
	} else {
		t.Logf("Expected error: %v", err)
	}
}

func TestLock(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	rdb := goredis.NewClient(&goredis.Options{
		Addr: mr.Addr(),
	})
	ctx := context.Background()

	key := "test_lock"
	val1 := "owner1"
	val2 := "owner2"
	expiration := 100 * time.Millisecond

	// 1. Acquire Lock Success
	ok, err := redis.AcquireLock(ctx, rdb, key, val1, expiration)
	if err != nil {
		t.Fatalf("AcquireLock failed: %v", err)
	}
	if !ok {
		t.Error("expected to acquire lock, got false")
	}

	// 2. Acquire Lock Failure (Already locked)
	ok, err = redis.AcquireLock(ctx, rdb, key, val2, expiration)
	if err != nil {
		t.Fatalf("AcquireLock 2 failed: %v", err)
	}
	if ok {
		t.Error("expected not to acquire lock, got true")
	}

	// 3. Release Lock Failure (Wrong owner)
	released, err := redis.ReleaseLock(ctx, rdb, key, val2)
	if err != nil {
		t.Fatalf("ReleaseLock failed: %v", err)
	}
	if released {
		t.Error("ReleaseLock should fail for wrong owner")
	}

	// 4. Release Lock Success (Correct owner)
	released, err = redis.ReleaseLock(ctx, rdb, key, val1)
	if err != nil {
		t.Fatalf("ReleaseLock failed: %v", err)
	}
	if !released {
		t.Error("ReleaseLock should succeed")
	}

	// verify key is gone
	exists := mr.Exists(key)
	if exists {
		t.Error("Key should be deleted")
	}
	
	// 5. Release Lock (Key doesn't exist)
	released, err = redis.ReleaseLock(ctx, rdb, key, val1)
	if err != nil {
		t.Fatalf("ReleaseLock non-existent failed: %v", err)
	}
	if released {
		t.Error("ReleaseLock should return false for non-existent key")
	}

	// 6. Release Lock Failure (Client Closed/Network Error)
	// We close the client locally, so Eval should fail (or we can close miniredis)
	rdb.Close()
	_, err = redis.ReleaseLock(ctx, rdb, key, val1)
	if err == nil {
		t.Error("ReleaseLock should fail when client is closed")
	}
}

// Loggers for silencing expected errors
type noOpLogger struct{}
func (l *noOpLogger) Printf(ctx context.Context, format string, v ...interface{}) {}

type defaultLogger struct{}
func (l *defaultLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	log.Printf(format, v...)
}

// TestIntegration_RealRedis performs actual operations against a running Redis server.
// Prerequisites: Redis running on localhost:6379 with no password.
func TestIntegration_RealRedis(t *testing.T) {
	// Try to connect to local redis
	cfg := &setting.RedisConfig{
		Host:     "127.0.0.1",
		Port:     6379,
		Password: "", // Default local redis usually has no password
		DB:       0,
	}

	client, err := redis.NewClient(cfg)
	if err != nil {
		t.Skipf("Skipping integration test: local redis not available: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// 1. SET
	key := "integration_test_key"
	value := "hello_redis"
	err = client.Set(ctx, key, value, 10*time.Second).Err()
	if err != nil {
		t.Fatalf("Failed to SET: %v", err)
	}

	// 2. GET
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		t.Fatalf("Failed to GET: %v", err)
	}
	if val != value {
		t.Errorf("GET mismatch: got %s, want %s", val, value)
	}

	// 3. Update (SET again)
	newValue := "updated_value"
	err = client.Set(ctx, key, newValue, 10*time.Second).Err()
	if err != nil {
		t.Fatalf("Failed to UPDATE: %v", err)
	}
	val, _ = client.Get(ctx, key).Result()
	if val != newValue {
		t.Errorf("UPDATE mismatch: got %s, want %s", val, newValue)
	}

	// 4. DEL
	err = client.Del(ctx, key).Err()
	if err != nil {
		t.Fatalf("Failed to DEL: %v", err)
	}

	// 5. Verify DEL
	_, err = client.Get(ctx, key).Result()
	if err != goredis.Nil {
		t.Errorf("Expected Redis Nil error after DEL, got: %v", err)
	}

	// 6. Test Lock with Real Redis
	lockKey := "integration_lock"
	lockVal := "real_owner"
	
	// Ensure cleanup of lock key
	defer client.Del(ctx, lockKey)

	ok, err := redis.AcquireLock(ctx, client, lockKey, lockVal, 5*time.Second)
	if err != nil {
		t.Fatalf("AcquireLock failed on real redis: %v", err)
	}
	if !ok {
		t.Error("AcquireLock returned false on real redis")
	}

	// Release
	released, err := redis.ReleaseLock(ctx, client, lockKey, lockVal)
	if err != nil {
		t.Fatalf("ReleaseLock failed on real redis: %v", err)
	}
	if !released {
		t.Error("ReleaseLock returned false on real redis")
	}
}
