package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/VishalHilal/e-commerce-api/internal/json"
)

type RateLimiter struct {
	clients map[string]*ClientLimiter
	mutex   sync.RWMutex
	rate    int
	window  time.Duration
}

type ClientLimiter struct {
	requests  int
	lastReset time.Time
	mutex     sync.Mutex
}

func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*ClientLimiter),
		rate:    rate,
		window:  window,
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := getClientIP(r)

		rl.mutex.RLock()
		client, exists := rl.clients[clientIP]
		rl.mutex.RUnlock()

		if !exists {
			rl.mutex.Lock()
			client = &ClientLimiter{
				lastReset: time.Now(),
			}
			rl.clients[clientIP] = client
			rl.mutex.Unlock()
		}

		client.mutex.Lock()

		// Reset counter if window has passed
		if time.Since(client.lastReset) > rl.window {
			client.requests = 0
			client.lastReset = time.Now()
		}

		client.requests++

		if client.requests > rl.rate {
			client.mutex.Unlock()
			json.WriteError(w, http.StatusTooManyRequests, "Rate limit exceeded")
			return
		}

		client.mutex.Unlock()
		next.ServeHTTP(w, r)
	})
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP if multiple are present
		if idx := strings.Index(xff, ","); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	// Fall back to RemoteAddr
	if idx := strings.LastIndex(r.RemoteAddr, ":"); idx != -1 {
		return r.RemoteAddr[:idx]
	}
	return r.RemoteAddr
}
