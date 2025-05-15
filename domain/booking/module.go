package booking

import "go.uber.org/fx"

// Module provides booking dependencies
var Module = fx.Module("booking",
	fx.Options(
		fx.Provide(
			NewRepository,
			NewService,
			NewController,
			NewRoute,
		),
		fx.Invoke(RegisterRoute),
	),
)
