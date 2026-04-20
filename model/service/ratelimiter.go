package service

import (
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	rate          *rate.Limiter
	lastEventNano atomic.Int64
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	lim := rate.NewLimiter(r, b)
	return &RateLimiter{rate: lim}
}

func (rm *RateLimiter) Allow() bool {
	rm.lastEventNano.Store(time.Now().UnixNano())
	return rm.rate.Allow()
}

func (rm *RateLimiter) LastEvent() time.Time {
	return time.Unix(0, rm.lastEventNano.Load())
}
