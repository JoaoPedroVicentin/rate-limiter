package client_entity

import (
	"context"
	"time"
)

type Client interface {
	Incr(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
}
