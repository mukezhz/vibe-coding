package domain

import (
	"clean-architecture/domain/booking"
	"clean-architecture/domain/organization"
	"clean-architecture/domain/todo"
	"clean-architecture/domain/user"

	"go.uber.org/fx"
)

var Module = fx.Module("domain",
	fx.Options(
		user.Module,
		todo.Module,
		organization.Module,
		booking.Module,
	),
)
