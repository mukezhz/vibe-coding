package todo

import (
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/infrastructure"
)

// Route struct
type Route struct {
	logger     framework.Logger
	handler    infrastructure.Router
	controller *Controller
}

// NewRoute creates a new route
func NewRoute(
	logger framework.Logger,
	handler infrastructure.Router,
	controller *Controller,
) *Route {
	return &Route{
		handler:    handler,
		logger:     logger,
		controller: controller,
	}
}

// RegisterRoute sets up todo routes
func RegisterRoute(r *Route) {
	r.logger.Info("Setting up todo routes")

	api := r.handler.Group("/api/todos")

	// Todo routes based on the .bru files
	api.POST("", r.controller.CreateTodo)
	api.GET("", r.controller.FetchTodoWithPagination)
	api.GET("/:id", r.controller.GetTodoByID)
	api.PUT("/:id", r.controller.UpdateTodo)
}
