package booking

import (
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/infrastructure"
)

// Route structure for booking system
type Route struct {
	logger     framework.Logger
	handler    infrastructure.Router
	controller *Controller
}

// NewRoute initializes booking routes
func NewRoute(
	logger framework.Logger,
	handler infrastructure.Router,
	controller *Controller,
) *Route {
	return &Route{
		logger:     logger,
		handler:    handler,
		controller: controller,
	}
}

// RegisterRoute configures booking API endpoints
func RegisterRoute(r *Route) {
	r.logger.Info("Setting up booking routes")

	// Group all API routes under /api
	api := r.handler.Group("/api")

	// Resource endpoints
	resources := api.Group("/resources")
	{
		resources.POST("", r.controller.CreateResource)
		resources.GET("", r.controller.ListResources)
		resources.GET("/:id", r.controller.GetResourceByID)
		resources.PUT("/:id", r.controller.UpdateResource)
		resources.DELETE("/:id", r.controller.DeleteResource)

		// Resource availability endpoints
		resources.GET("/:id/availability", r.controller.CheckResourceAvailability)
		resources.POST("/:id/availability", r.controller.CreateAvailability)
		resources.GET("/:id/availabilities", r.controller.ListResourceAvailabilities)
	}

	// Availability endpoints for checking multiple resources
	api.GET("/availability", r.controller.CheckMultipleResourcesAvailability)

	// Booking endpoints
	bookings := api.Group("/bookings")
	{
		bookings.POST("", r.controller.CreateBooking)
		bookings.GET("", r.controller.ListBookings)
		bookings.GET("/:id", r.controller.GetBookingByID)
		bookings.PUT("/:id", r.controller.UpdateBooking)
		bookings.DELETE("/:id", r.controller.CancelBooking)
	}

	// User bookings endpoint
	api.GET("/users/:id/bookings", r.controller.ListUserBookings)
}
