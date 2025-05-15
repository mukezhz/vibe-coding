package booking

import (
	"clean-architecture/domain/models"
	"clean-architecture/pkg/types"
	"time"
)

// ResourceCreateDTO for creating a new resource
type ResourceCreateDTO struct {
	Name        string                 `json:"name" binding:"required"`
	Description string                 `json:"description"`
	Type        string                 `json:"type" binding:"required"`
	Capacity    int                    `json:"capacity"`
	Location    string                 `json:"location"`
	Attributes  map[string]interface{} `json:"attributes"`
}

// ResourceResponseDTO for resource responses
type ResourceResponseDTO struct {
	UUID        string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        string                 `json:"type"`
	Capacity    int                    `json:"capacity"`
	Location    string                 `json:"location"`
	Attributes  map[string]interface{} `json:"attributes"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// ResourceUpdateDTO for updating a resource
type ResourceUpdateDTO struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        string                 `json:"type"`
	Capacity    int                    `json:"capacity"`
	Location    string                 `json:"location"`
	Attributes  map[string]interface{} `json:"attributes"`
}

// AvailabilityCreateDTO for creating availability
type AvailabilityCreateDTO struct {
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required"`
	IsRecurring bool      `json:"is_recurring"`
	RecurRule   string    `json:"recur_rule"`
}

// AvailabilityResponseDTO for availability responses
type AvailabilityResponseDTO struct {
	UUID        string    `json:"id"`
	ResourceID  string    `json:"resource_id"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	IsRecurring bool      `json:"is_recurring"`
	RecurRule   string    `json:"recur_rule"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AvailabilityUpdateDTO for updating availability
type AvailabilityUpdateDTO struct {
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	IsRecurring bool      `json:"is_recurring"`
	RecurRule   string    `json:"recur_rule"`
}

// AvailabilityCheckDTO for checking availability
type AvailabilityCheckDTO struct {
	StartTime time.Time `json:"start_time" binding:"required" form:"start"`
	EndTime   time.Time `json:"end_time" binding:"required" form:"end"`
}

// AvailabilityCheckResponseDTO for availability check responses
type AvailabilityCheckResponseDTO struct {
	Available bool `json:"available"`
}

// BookingCreateDTO for creating a booking
type BookingCreateDTO struct {
	ResourceID types.BinaryUUID `json:"resource_id" binding:"required"`
	StartTime  time.Time        `json:"start_time" binding:"required"`
	EndTime    time.Time        `json:"end_time" binding:"required"`
	Notes      string           `json:"notes"`
	Reference  string           `json:"reference"`
}

// BookingResponseDTO for booking responses
type BookingResponseDTO struct {
	UUID       string    `json:"id"`
	ResourceID string    `json:"resource_id"`
	UserID     string    `json:"user_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Status     string    `json:"status"`
	Notes      string    `json:"notes"`
	Reference  string    `json:"reference"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// BookingUpdateDTO for updating a booking
type BookingUpdateDTO struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
	Notes     string    `json:"notes"`
	Reference string    `json:"reference"`
}

// ResourceQueryParams for filtering resources
type ResourceQueryParams struct {
	Type     string `form:"type"`
	Location string `form:"location"`
	Capacity int    `form:"capacity"`
	Page     int    `form:"page,default=1"`
	Limit    int    `form:"limit,default=10"`
}

// BookingQueryParams for filtering bookings
type BookingQueryParams struct {
	ResourceID string    `form:"resource_id"`
	UserID     string    `form:"user_id"`
	StartTime  time.Time `form:"start_time"`
	EndTime    time.Time `form:"end_time"`
	Status     string    `form:"status"`
	Page       int       `form:"page,default=1"`
	Limit      int       `form:"limit,default=10"`
}

// ResourceToDTO converts a Resource model to ResourceResponseDTO
func ResourceToDTO(resource *models.Resource) ResourceResponseDTO {
	var attributes map[string]interface{}
	if resource.Attributes != nil {
		_ = resource.Attributes.UnmarshalJSON([]byte(resource.Attributes.String()))
	}

	return ResourceResponseDTO{
		UUID:        resource.UUID.String(),
		Name:        resource.Name,
		Description: resource.Description,
		Type:        resource.Type,
		Capacity:    resource.Capacity,
		Location:    resource.Location,
		Attributes:  attributes,
		CreatedAt:   resource.CreatedAt,
		UpdatedAt:   resource.UpdatedAt,
	}
}

// AvailabilityToDTO converts an Availability model to AvailabilityResponseDTO
func AvailabilityToDTO(availability *models.Availability) AvailabilityResponseDTO {
	return AvailabilityResponseDTO{
		UUID:        availability.UUID.String(),
		ResourceID:  availability.ResourceID.String(),
		StartTime:   availability.StartTime,
		EndTime:     availability.EndTime,
		IsRecurring: availability.IsRecurring,
		RecurRule:   availability.RecurRule,
		CreatedAt:   availability.CreatedAt,
		UpdatedAt:   availability.UpdatedAt,
	}
}

// BookingToDTO converts a Booking model to BookingResponseDTO
func BookingToDTO(booking *models.Booking) BookingResponseDTO {
	return BookingResponseDTO{
		UUID:       booking.UUID.String(),
		ResourceID: booking.ResourceID.String(),
		UserID:     booking.UserID.String(),
		StartTime:  booking.StartTime,
		EndTime:    booking.EndTime,
		Status:     booking.Status,
		Notes:      booking.Notes,
		Reference:  booking.Reference,
		CreatedAt:  booking.CreatedAt,
		UpdatedAt:  booking.UpdatedAt,
	}
}
