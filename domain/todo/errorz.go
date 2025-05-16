package todo

import (
	"clean-architecture/pkg/errorz"
)

var (
	ErrInvalidTodoID     = errorz.ErrBadRequest.JoinError("Invalid Todo ID")
	ErrTodoNotFound      = errorz.ErrNotFound.JoinError("Todo not found")
	ErrTodoTitleRequired = errorz.ErrBadRequest.JoinError("Todo title is required")
)
