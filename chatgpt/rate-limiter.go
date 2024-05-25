package chatgpt

import (
	"sync"
	"time"
)

// LeakyBucket represents the leaky bucket for each ID
type LeakyBucket struct {
	capacity    int
	remaining   int
	lastUpdated time.Time
	ratePerSec  float64
	mu          sync.Mutex
}

// RateLimiter stores the leaky buckets for different IDs
type RateLimiter[T comparable] struct {
	buckets    map[T]*LeakyBucket
	mu         sync.Mutex
	capacity   int
	ratePerSec float64
}

// NewRateLimiter creates a new RateLimiter
func NewRateLimiter[T comparable](capacity int, maxRequests int, interval time.Duration) *RateLimiter[T] {
	return &RateLimiter[T]{
		buckets:    make(map[T]*LeakyBucket),
		capacity:   capacity,
		ratePerSec: float64(maxRequests) / interval.Seconds(),
	}
}

// Allow checks if a request from the given ID can go through or should be blocked
func (rl *RateLimiter[T]) Allow(id T) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, exists := rl.buckets[id]
	if !exists {
		bucket = &LeakyBucket{
			capacity:    rl.capacity,
			remaining:   rl.capacity,
			lastUpdated: time.Now(),
			ratePerSec:  rl.ratePerSec,
		}
		rl.buckets[id] = bucket
	}

	return bucket.allowRequest()
}

// allowRequest checks and updates the bucket state to determine if a request can go through
func (lb *LeakyBucket) allowRequest() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(lb.lastUpdated).Seconds()

	// Update the remaining tokens considering the elapsed time
	lb.remaining += int(elapsed * lb.ratePerSec)
	if lb.remaining > lb.capacity {
		lb.remaining = lb.capacity
	}
	lb.lastUpdated = now

	if lb.remaining > 0 {
		lb.remaining--
		return true
	}

	return false
}
