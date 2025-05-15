package booking

import "clean-architecture/pkg/errorz"

// Domain-specific error codes for the booking system
const (
	ErrCodeResourceNotFound     = "RESOURCE_NOT_FOUND"
	ErrCodeBookingNotFound      = "BOOKING_NOT_FOUND"
	ErrCodeAvailabilityNotFound = "AVAILABILITY_NOT_FOUND"
	ErrCodeResourceNotAvailable = "RESOURCE_NOT_AVAILABLE"
	ErrCodeInvalidTimeRange     = "INVALID_TIME_RANGE"
	ErrCodeBookingOverlap       = "BOOKING_OVERLAP"
	ErrCodeInvalidBookingStatus = "INVALID_BOOKING_STATUS"
	ErrCodePastDateBooking      = "PAST_DATE_BOOKING"
	ErrCodeExceedsMaxDuration   = "EXCEEDS_MAX_DURATION"
	ErrCodeInsufficientLeadTime = "INSUFFICIENT_LEAD_TIME"
)

var (
	// ErrResourceNotFound is returned when a resource is not found
	ErrResourceNotFound = errorz.ErrNotFound.JoinError("resource not found")

	// ErrBookingNotFound is returned when a booking is not found
	ErrBookingNotFound = errorz.ErrNotFound.JoinError("booking not found")

	// ErrAvailabilityNotFound is returned when availability is not found
	ErrAvailabilityNotFound = errorz.ErrNotFound.JoinError("availability not found")

	// ErrResourceNotAvailable is returned when a resource is not available for the requested time
	ErrResourceNotAvailable = errorz.ErrConflict.JoinError("resource not available for the requested time period")

	// ErrInvalidTimeRange is returned when an invalid time range is provided
	ErrInvalidTimeRange = errorz.ErrBadRequest.JoinError("invalid time range")

	// ErrBookingOverlap is returned when a booking overlaps with existing bookings
	ErrBookingOverlap = errorz.ErrConflict.JoinError("booking overlaps with existing bookings")

	// ErrInvalidBookingStatus is returned when an invalid booking status is provided
	ErrInvalidBookingStatus = errorz.ErrBadRequest.JoinError("invalid booking status")

	// ErrPastDateBooking is returned when attempting to book in the past
	ErrPastDateBooking = errorz.ErrBadRequest.JoinError("cannot book in the past")

	// ErrExceedsMaxDuration is returned when a booking exceeds the maximum allowed duration
	ErrExceedsMaxDuration = errorz.ErrBadRequest.JoinError("booking exceeds maximum allowed duration")

	// ErrInsufficientLeadTime is returned when a booking doesn't meet the minimum lead time requirement
	ErrInsufficientLeadTime = errorz.ErrBadRequest.JoinError("booking does not meet minimum lead time requirement")
)
