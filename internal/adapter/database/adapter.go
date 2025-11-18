package database

import (
	"example-otelgo-miniservice/internal/config"
	"example-otelgo-miniservice/internal/domain"

	"go.uber.org/fx"
)

var Module = fx.Module("database_adapter",
	fx.Provide(func() domain.MongoOptions {
		return domain.MongoOptions{
			URL:            config.GetEnv("MONGO_URL", "mongodb://localhost:27017"),
			DBName:         config.GetEnv("MONGO_DB_NAME", "example-otelgo-miniservice"),
			ConnTimeoutSec: config.GetEnvAsInt("MONGO_CONN_TIMEOUT", 10),
			OpTimeoutSec:   config.GetEnvAsInt("MONGO_OP_TIMEOUT", 10),
			TxnTimeoutSec:  config.GetEnvAsInt("MONGO_TXN_TIMEOUT", 10),
		}
	}),
	fx.Provide(NewMongoDB),
	fx.Provide(NewRepository[domain.UserRepo]),
	fx.Provide(NewRepository[domain.RoleRepo]),
)
