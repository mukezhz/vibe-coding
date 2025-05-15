package organization

import (
	"go.uber.org/fx"
)

// Module exports organization dependencies
var Module = fx.Module("organization",
	fx.Provide(
		NewRepository,
		NewService,
		NewController,
		NewRoute,
	),
	fx.Invoke(RegisterRoutes),
)
