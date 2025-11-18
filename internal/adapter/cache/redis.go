package cache

import (
	"context"
	"encoding/json"
	"errors"
	"example-otelgo-miniservice/internal/config"
	"example-otelgo-miniservice/internal/port"
	"example-otelgo-miniservice/internal/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedis() (port.Redis, error) {
	var client *redis.Client

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := utils.Retry(ctx, utils.RetryConfig{
		MaxAttempts: 5,
		Interval:    3 * time.Second,
		Backoff:     true,
	}, func(attempt int) error {

		client = redis.NewClient(&redis.Options{
			Addr:     config.GetEnv("REDIS_ADDR", "localhost:6379"),
			Username: config.GetEnv("REDIS_USER", "default"),
			Password: config.GetEnv("REDIS_PASS", ""),
			DB:       config.GetEnvAsInt("REDIS_DB", 0),
		})

		return client.Ping(ctx).Err()
	}); err != nil {
		return nil, err
	}

	return &Redis{client: client}, nil
}

func (r *Redis) Client() *redis.Client {
	return r.client
}

func (r *Redis) Close(ctx context.Context) error {
	return r.client.Close()
}

func (r *Redis) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil && !r.IsNil(err) {
		return nil, err
	}
	return data, err
}

func (r *Redis) Del(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return err
}

func (r *Redis) SetWithTTL(ctx context.Context, key string, value string, ttl time.Duration) error {
	err := r.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}
	return err
}

func (r *Redis) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func (r *Redis) Incr(ctx context.Context, key string) (int64, error) {
	result, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return result, err
}

func (r *Redis) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return ttl, err
}

func (r *Redis) IsNil(err error) bool {
	return errors.Is(err, redis.Nil)
}
