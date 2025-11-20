package rate_limiter

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/JoaoPedroVicentin/rate-limiter/internal/entity/client_entity"
)

type RateLimiter struct {
	client        client_entity.Client
	limitIp       int
	limitApiKey   int
	window        time.Duration
	blockDuration time.Duration
	context       context.Context
}

func NewRateLimiter(client client_entity.Client, limitIp int, limitApiKey int, window time.Duration, blockDuration time.Duration) *RateLimiter {
	return &RateLimiter{
		client:        client,
		limitIp:       limitIp,
		limitApiKey:   limitApiKey,
		window:        window,
		blockDuration: blockDuration,
		context:       context.Background(),
	}
}

func (rl *RateLimiter) Allow(ip string, apiKey string, apiKeyHeader string) (bool, error) {
	var limit int
	var key string

	if apiKey == apiKeyHeader {
		limit = rl.limitApiKey
		key = "apiKey:" + apiKeyHeader
	} else {
		limit = rl.limitIp
		key = "ip:" + ip
	}

	incr, err := rl.client.Incr(rl.context, key)
	if err != nil {
		return false, err
	}
	err = rl.client.Expire(rl.context, key, rl.window)
	if err != nil {
		return false, err
	}

	if incr > int64(limit) {
		_ = rl.client.Expire(rl.context, key, rl.blockDuration)
		return false, nil
	}

	return true, nil
}

func RateLimiterMiddleware(rl *RateLimiter, next http.Handler, apiKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		apiKeyHeader := r.Header.Get("api-key")

		clientIp, _, _ := net.SplitHostPort(r.RemoteAddr)

		allowed, err := rl.Allow(clientIp, apiKey, apiKeyHeader)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		if !allowed {
			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
