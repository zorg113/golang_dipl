package service

import (
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	rate      *rate.Limiter
	LastEvent time.Time
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	lim := rate.NewLimiter(r, b)
	return &RateLimiter{rate: lim}
}

func (rm *RateLimiter) Allow() bool {
	rm.LastEvent = time.Now()
	allow := rm.rate.Allow()
	return allow
}
