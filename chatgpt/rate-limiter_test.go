package chatgpt

import (
	"testing"
	"time"
)

func TestAllow(t *testing.T) {
	capacity := 10
	ratePerMinute := 10
	rl := NewRateLimiter[string](capacity, ratePerMinute, time.Minute)
	id := "user1"

	// Test allowing requests within limit
	for i := 0; i < 10; i++ {
		if !rl.Allow(id) {
			t.Errorf("Request %d for %s was blocked but should have been allowed", i+1, id)
		}
	}

	// Test blocking requests exceeding limit
	for i := 0; i < 5; i++ {
		if rl.Allow(id) {
			t.Errorf("Request %d for %s was allowed but should have been blocked", i+11, id)
		}
	}

	// Wait for enough time for the rate limiter to allow one more request
	time.Sleep(6 * time.Second)

	// Test allowing request after waiting for the bucket to refill
	if !rl.Allow(id) {
		t.Errorf("Request %d for %s was blocked but should have been allowed", 17, id)
	}

	// Test blocking request because we didn't wait long enough
	if rl.Allow(id) {
		t.Errorf("Request %d for %s was allowed but should have been blocked", 18, id)
	}

	// Wait for enough time for the rate limiter to allow one more request again
	time.Sleep(6 * time.Second)

	// Test allowing request after waiting for the bucket to refill again
	if !rl.Allow(id) {
		t.Errorf("Request %d for %s was blocked but should have been allowed", 19, id)
	}
}

func TestAllowMultipleIDs(t *testing.T) {
	capacity := 10
	ratePerMinute := 5
	rl := NewRateLimiter[string](capacity, ratePerMinute, time.Minute)

	id1 := "user1"
	id2 := "user2"

	// Test allowing requests for multiple IDs
	for i := 0; i < 10; i++ {
		if !rl.Allow(id1) {
			t.Errorf("Request %d for %s was blocked but should have been allowed", i+1, id1)
		}
		if !rl.Allow(id2) {
			t.Errorf("Request %d for %s was blocked but should have been allowed", i+1, id2)
		}
	}

	// Test blocking requests exceeding limit for both IDs
	for i := 0; i < 5; i++ {
		if rl.Allow(id1) {
			t.Errorf("Request %d for %s was allowed but should have been blocked", i+11, id1)
		}
		if rl.Allow(id2) {
			t.Errorf("Request %d for %s was allowed but should have been blocked", i+11, id2)
		}
	}
}

func TestRateLimiterConcurrency(t *testing.T) {
	capacity := 10
	ratePerMinute := 5
	rl := NewRateLimiter[string](capacity, ratePerMinute, time.Minute)

	id := "user1"
	numGoroutines := 100
	allowed := make(chan bool, numGoroutines)

	// Spawn multiple goroutines to test concurrency
	for i := 0; i < numGoroutines; i++ {
		go func() {
			allowed <- rl.Allow(id)
		}()
	}

	// Collect results
	allowedCount := 0
	blockedCount := 0
	for i := 0; i < numGoroutines; i++ {
		if <-allowed {
			allowedCount++
		} else {
			blockedCount++
		}
	}

	// Verify results
	if allowedCount > capacity {
		t.Errorf("Allowed count %d exceeds capacity %d", allowedCount, capacity)
	}

	if blockedCount != numGoroutines-allowedCount {
		t.Errorf("Blocked count %d does not match expected %d", blockedCount, numGoroutines-allowedCount)
	}
}
