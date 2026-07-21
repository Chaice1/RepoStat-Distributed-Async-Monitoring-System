package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

type RedisRepo interface {
	CheckRateLimit(context.Context, string, int, float64) (bool, error)
}

type rateLimiter struct {
	redis   RedisRepo
	log     *slog.Logger
	mu      *sync.Mutex
	storage map[string]*rate.Limiter
	burst   int
	rps     float64
}

func NewRateLimiter(log *slog.Logger, redis RedisRepo, burst int, rps float64) *rateLimiter {
	return &rateLimiter{
		log:     log,
		redis:   redis,
		storage: make(map[string]*rate.Limiter, 0),
		burst:   burst,
		rps:     rps,
	}
}

func (rl *rateLimiter) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		allowed, err := rl.redis.CheckRateLimit(context.Background(), ip, rl.burst, rl.rps)

		if err != nil {

			rl.log.Warn("redis is down", "error", err)

			if !rl.checkRateLimitLocal(ip) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				if err := json.NewEncoder(w).Encode(map[string]string{"error": "rate limit exceeded"}); err != nil {
					rl.log.Error("failed to write message to jsonbody", "error", err)
				}
				return
			}

		} else if !allowed {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			if err := json.NewEncoder(w).Encode(map[string]string{"error": "rate limit exceeded"}); err != nil {
				rl.log.Error("failed to write message to jsonbody", "error", err)
			}
			return
		}

		next.ServeHTTP(w, r)

	})
}

func (rl *rateLimiter) checkRateLimitLocal(ip string) bool {

	rl.mu.Lock()
	defer rl.mu.Unlock()

	l, ok := rl.storage[ip]

	if !ok {
		l = rate.NewLimiter(rate.Limit(rl.rps), rl.burst)
		rl.storage[ip] = l
	}

	return l.Allow()
}
