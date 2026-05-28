package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"github.com/brandon200217/NOTIFY/utilities"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
	rps      rate.Limit // requests por segundo permitidos
	burst    int        // tamaño (ráfaga máxima)
}

func NewRateLimiter(rps float64, burst int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rps:      rate.Limit(rps),
		burst:    burst,
	}
}

func (rl *RateLimiter) getLimiter(token string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[token]
	if !exists {
		limiter = rate.NewLimiter(rl.rps, rl.burst)
		rl.limiters[token] = limiter
	}
	return limiter
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := utilities.ExtractBearerToken(r)
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}

		limiter := rl.getLimiter(token)
		if !limiter.Allow() {
			slog.Warn("rate limit excedido",
				"remote_addr", r.RemoteAddr)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":  "demasiados requests, intentá más tarde",
				"status": http.StatusTooManyRequests,
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
