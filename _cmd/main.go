package main

import (
	"context"
	"example-otelgo-miniservice/internal/adapter/cache"
	"example-otelgo-miniservice/internal/adapter/database"
	"example-otelgo-miniservice/internal/api"
	"example-otelgo-miniservice/internal/application"
	"example-otelgo-miniservice/internal/config"
	"example-otelgo-miniservice/internal/port"
	"log"
	"time"

	"github.com/Maximumsoft-Co-LTD/otelgo/eto"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Invoke(config.LoadEnv),
		database.Module,
		cache.Module,
		application.Module,
		api.Module,
		fx.Invoke(RunApp),
	)

	app.Run()
}

func RunApp(
	lc fx.Lifecycle,
	api port.Router,
	mongodb port.MongoDB,
	redis port.Redis,
) {
	var shutdown func(ctx context.Context) error
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			var err error
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			shutdown, err = eto.Init(ctx, eto.Config{
				ServiceName:   config.GetEnv("SERVICE_NAME", "example-otelgo-miniservice"),
				Environment:   config.GetEnv("SERVICE_MODE", "development"),
				OtelEndpoint:  config.GetEnv("OTEL_ENDPOINT", "otel-collector:4317"),
				EnableMetrics: config.GetEnvAsBool("ENABLE_METRICS", true),
			})
			if err != nil {
				log.Fatalf("eto init error: %v", err)
			}

			go func() {
				if err := api.Start(); err != nil {
					panic(err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := mongodb.Close(ctx); err != nil {
				return err
			}

			if err := redis.Close(ctx); err != nil {
				return err
			}

			if err := shutdown(ctx); err != nil {
				return err
			}

			return nil
		},
	})
}
