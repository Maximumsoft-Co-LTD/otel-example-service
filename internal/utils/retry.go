package utils

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
)

type RetryConfig struct {
	MaxAttempts int
	Interval    time.Duration
	Backoff     bool
}

// Retry - รองรับ attempt แบบระบุจำนวนรอบได้
func Retry(ctx context.Context, cfg RetryConfig, fn func(attempt int) error) error {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 3
	}
	if cfg.Interval <= 0 {
		cfg.Interval = 500 * time.Millisecond
	}

	attempt := 0
	var strategy backoff.BackOff

	if cfg.Backoff {
		exp := backoff.NewExponentialBackOff()
		exp.InitialInterval = cfg.Interval
		exp.MaxElapsedTime = time.Duration(cfg.MaxAttempts) * cfg.Interval
		strategy = exp
	} else {
		constant := backoff.NewConstantBackOff(cfg.Interval)
		strategy = backoff.WithMaxRetries(constant, uint64(cfg.MaxAttempts))
	}

	bo := backoff.WithContext(strategy, ctx)

	return backoff.Retry(func() error {
		attempt++
		return fn(attempt)
	}, bo)
}
