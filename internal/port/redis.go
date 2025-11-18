package port

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis interface {
	Client() *redis.Client
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Del(ctx context.Context, key string) error
	SetWithTTL(ctx context.Context, key string, value string, ttl time.Duration) error
	Exists(ctx context.Context, key string) (bool, error)
	Incr(ctx context.Context, key string) (int64, error)
	TTL(ctx context.Context, key string) (time.Duration, error)
	IsNil(err error) bool
	Close(ctx context.Context) error
}
