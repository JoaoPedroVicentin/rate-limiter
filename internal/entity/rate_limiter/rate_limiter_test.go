package rate_limiter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockClient struct {
	calls   map[string]int64
	expires map[string]time.Time
}

func newMockClient() *mockClient {
	return &mockClient{
		calls:   make(map[string]int64),
		expires: make(map[string]time.Time),
	}
}

func (m *mockClient) Incr(ctx context.Context, key string) (int64, error) {
	if exp, ok := m.expires[key]; ok && time.Now().After(exp) {
		m.calls[key] = 0
	}
	m.calls[key]++
	return m.calls[key], nil
}

func (m *mockClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	m.expires[key] = time.Now().Add(expiration)
	return nil
}

func TestRateLimiter_Allow_IPLimit(t *testing.T) {
	mock := newMockClient()
	rl := NewRateLimiter(mock, 2, 0, time.Second, 2*time.Second)

	ip := "127.0.0.1"
	apiKey := "12345"
	apiKeyHeader := ""

	// 1st request: allowed
	allowed, err := rl.Allow(ip, apiKey, apiKeyHeader)
	assert.NoError(t, err)
	assert.True(t, allowed)

	// 2nd request: allowed
	allowed, err = rl.Allow(ip, apiKey, apiKeyHeader)
	assert.NoError(t, err)
	assert.True(t, allowed)

	// 3rd request: blocked
	allowed, err = rl.Allow(ip, apiKey, apiKeyHeader)
	assert.NoError(t, err)
	assert.False(t, allowed)

	// Wait for block duration to expire
	time.Sleep(2100 * time.Millisecond)
	allowed, err = rl.Allow(ip, apiKey, apiKeyHeader)
	assert.NoError(t, err)
	assert.True(t, allowed)
}

func TestRateLimiter_Allow_ApiKeyLimit(t *testing.T) {
	mock := newMockClient()
	rl := NewRateLimiter(mock, 0, 2, time.Second, 2*time.Second)

	ip := "127.0.0.1"
	apiKey := "12345"
	apiKeyHeader := "12345"

	// 1st request: allowed
	allowed, err := rl.Allow(ip, apiKey, apiKeyHeader)
	assert.NoError(t, err)
	assert.True(t, allowed)

	// 2nd request: allowed
	allowed, err = rl.Allow(ip, apiKey, apiKeyHeader)
	assert.NoError(t, err)
	assert.True(t, allowed)

	// 3rd request: blocked
	allowed, err = rl.Allow(ip, apiKey, apiKeyHeader)
	assert.NoError(t, err)
	assert.False(t, allowed)

	// Wait for block duration to expire
	time.Sleep(2100 * time.Millisecond)
	allowed, err = rl.Allow(ip, apiKey, apiKeyHeader)
	assert.NoError(t, err)
	assert.True(t, allowed)
}

func TestRateLimiterMiddleware_Integration(t *testing.T) {
	mock := newMockClient()
	rl := NewRateLimiter(mock, 2, 0, time.Second, 2*time.Second)

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := RateLimiterMiddleware(rl, h, "12345")

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "127.0.0.1:1234"

	// 1st request: allowed
	rw := httptest.NewRecorder()
	mw.ServeHTTP(rw, req)
	assert.Equal(t, http.StatusOK, rw.Code)

	// 2nd request: allowed
	rw = httptest.NewRecorder()
	mw.ServeHTTP(rw, req)
	assert.Equal(t, http.StatusOK, rw.Code)

	// 3rd request: blocked
	rw = httptest.NewRecorder()
	mw.ServeHTTP(rw, req)
	assert.Equal(t, http.StatusTooManyRequests, rw.Code)

	// Wait for block duration to expire
	time.Sleep(2100 * time.Millisecond)

	// 4th request: allowed again
	rw = httptest.NewRecorder()
	mw.ServeHTTP(rw, req)
	assert.Equal(t, http.StatusOK, rw.Code)
}
