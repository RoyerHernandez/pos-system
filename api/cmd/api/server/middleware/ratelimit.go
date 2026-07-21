package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type visitor struct {
	count    int
	lastSeen time.Time
}

// RateLimiter tracks request counts per IP within a sliding window.
type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a rate limiter that allows limit requests per window per IP.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		limit:    limit,
		window:   window,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(rl.window)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.window {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Handler is a chi-compatible middleware that enforces the rate limit.
func (rl *RateLimiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		rl.mu.Lock()
		v, exists := rl.visitors[ip]
		if !exists {
			rl.visitors[ip] = &visitor{count: 1, lastSeen: time.Now()}
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if time.Since(v.lastSeen) > rl.window {
			v.count = 1
			v.lastSeen = time.Now()
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		v.count++
		v.lastSeen = time.Now()
		if v.count > rl.limit {
			rl.mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error":{"code":"RATE_LIMITED","message":"too many requests"}}`))
			return
		}
		rl.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
