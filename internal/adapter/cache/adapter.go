package cache

import (
	"example-otelgo-miniservice/internal/config"
	"example-otelgo-miniservice/internal/domain"

	"go.uber.org/fx"
)

var Module = fx.Module("cache_adapter",
	fx.Provide(func() domain.RedisOptions {
		return domain.RedisOptions{
			Addr:     config.GetEnv("REDIS_ADDR", "localhost:6379"),
			Password: config.GetEnv("REDIS_PASSWORD", ""),
			DB:       config.GetEnvAsInt("REDIS_DB", 0),
		}
	}),
	fx.Provide(NewRedis),
)
