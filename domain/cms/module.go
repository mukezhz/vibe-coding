package cms

import "go.uber.org/fx"

var Module = fx.Module("cms",
	fx.Options(
		fx.Provide(
			NewRepository,
			NewService,
			NewController,
			NewRoute,
		),
		fx.Invoke(Migrate, RegisterRoute),
	),
)
