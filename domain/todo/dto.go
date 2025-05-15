package todo

import (
	"clean-architecture/pkg/responses"
	"time"
)

// CreateTodoRequest DTO for creating a todo
type CreateTodoRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

// TodoResponse DTO for todo response
type TodoResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TodoListItem DTO for items in todo list
type TodoListItem struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// TodoListResponse DTO for paginated todo list
type TodoListResponse = responses.ListResponseType[TodoListItem]

// UpdateTodoRequest DTO for updating a todo
type UpdateTodoRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}
