package application

import (
	"example-otelgo-miniservice/internal/application/core"
	"example-otelgo-miniservice/internal/application/service"

	"go.uber.org/fx"
)

var Module = fx.Module("application",
	// SERVICE
	fx.Provide(service.NewJWE),
	// CORE
	fx.Provide(core.NewAuth),
)
