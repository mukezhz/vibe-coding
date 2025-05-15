package organization

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

// RegisterRoutes registers the organization routes
func RegisterRoutes(r *Route) {
	api := r.handler.Group("/api/organizations")
	api.POST("", r.controller.Create)
	api.GET("", r.controller.List)
	api.GET("/:id", r.controller.GetByID)
	api.PUT("/:id", r.controller.Update)
}
