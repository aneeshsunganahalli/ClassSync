package middlewares

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	mu sync.Mutex
	visitors map[string]int
	limit int
	resetTime time.Duration
}

func NewRateLimiter(limit int, resetTime time.Duration) *rateLimiter {
	rl := &rateLimiter{
		visitors: make(map[string]int),
		limit: limit,
		resetTime: resetTime,
	}

	go rl.resetVisitorCount() // Needs to be out main thread since we need to infinitely reset after resetTime
	return rl
}

func(rl *rateLimiter) resetVisitorCount() {
	for {
		time.Sleep(rl.resetTime)
		rl.mu.Lock()
		rl.visitors = make(map[string]int)
		rl.mu.Unlock()
	}
}

func (rl *rateLimiter) Middleware(next http.Handler) http.Handler {
		return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {

				rl.mu.Lock()
				defer rl.mu.Unlock()
				visitorIP := r.RemoteAddr
				fmt.Printf("%v made %v requests\n", visitorIP, rl.visitors[visitorIP])
				rl.visitors[visitorIP]++
				if rl.visitors[visitorIP] > rl.limit {
					http.Error(w, "Too many requests\n", http.StatusTooManyRequests)
					return
				}

				next.ServeHTTP(w,r)
		})
}