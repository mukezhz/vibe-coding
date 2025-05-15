package todo

import (
	"clean-architecture/pkg/errorz"
)

// Error codes specific to the todo domain
const (
	ErrInvalidTodoID     = "INVALID_TODO_ID"
	ErrTodoNotFound      = "TODO_NOT_FOUND"
	ErrTodoTitleRequired = "TODO_TITLE_REQUIRED"
)

// NewInvalidTodoIDError returns a new error for invalid todo ID
func NewInvalidTodoIDError() error {
	return errorz.ErrBadRequest.JoinError("Invalid Todo ID")
}

// NewTodoNotFoundError returns a new error when todo is not found
func NewTodoNotFoundError() error {
	return errorz.ErrNotFound.JoinError("Todo not found")
}

// NewTodoTitleRequiredError returns a new error when todo title is missing
func NewTodoTitleRequiredError() error {
	return errorz.ErrBadRequest.JoinError("Todo title is required")
}
