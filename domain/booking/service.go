package booking

import (
	"errors"
	"time"

	"clean-architecture/domain/models"
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service contains business logic for booking system
type Service struct {
	logger     framework.Logger
	repository Repository
}

// NewService creates a new booking service
func NewService(logger framework.Logger, repository Repository) *Service {
	return &Service{
		logger:     logger,
		repository: repository,
	}
}

// -------------- Resource Service Methods --------------

// CreateResource creates a new resource
func (s *Service) CreateResource(resource *models.Resource) error {
	s.logger.Info("[BookingService...CreateResource]")
	return s.repository.CreateResource(resource)
}

// GetResourceByID gets a resource by ID
func (s *Service) GetResourceByID(id types.BinaryUUID) (models.Resource, error) {
	s.logger.Info("[BookingService...GetResourceByID]")

	resource, err := s.repository.GetResourceByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resource, ErrResourceNotFound
		}
		return resource, err
	}

	return resource, nil
}

// UpdateResource updates a resource
func (s *Service) UpdateResource(id types.BinaryUUID, updateFn func(*models.Resource) error) error {
	s.logger.Info("[BookingService...UpdateResource]")

	// Get existing resource
	resource, err := s.repository.GetResourceByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrResourceNotFound
		}
		return err
	}

	// Apply updates via callback function
	if err := updateFn(&resource); err != nil {
		return err
	}

	// Save updated resource
	return s.repository.UpdateResource(&resource)
}

// DeleteResource deletes a resource
func (s *Service) DeleteResource(id types.BinaryUUID) error {
	s.logger.Info("[BookingService...DeleteResource]")

	// Check if resource exists
	_, err := s.repository.GetResourceByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrResourceNotFound
		}
		return err
	}

	// Delete resource
	return s.repository.DeleteResource(id)
}

// ListResources lists resources with pagination and filtering
func (s *Service) ListResources(page, limit int, filters map[string]interface{}) ([]models.Resource, int64, error) {
	s.logger.Info("[BookingService...ListResources]")
	return s.repository.ListResources(page, limit, filters)
}

// -------------- Availability Service Methods --------------

// CreateAvailability creates a new availability
func (s *Service) CreateAvailability(resourceID types.BinaryUUID, availability *models.Availability) error {
	s.logger.Info("[BookingService...CreateAvailability]")

	// Validate time range
	if availability.EndTime.Before(availability.StartTime) || availability.StartTime.Before(time.Now()) {
		return ErrInvalidTimeRange
	}

	// Check if resource exists
	_, err := s.repository.GetResourceByID(resourceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrResourceNotFound
		}
		return err
	}

	// Set resource ID
	availability.ResourceID = resourceID

	// Generate UUID if not provided
	if availability.UUID.String() == (types.BinaryUUID{}).String() {
		id, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		availability.UUID = types.BinaryUUID(id)
	}

	return s.repository.CreateAvailability(availability)
}

// GetAvailabilityByID gets an availability by ID
func (s *Service) GetAvailabilityByID(id types.BinaryUUID) (models.Availability, error) {
	s.logger.Info("[BookingService...GetAvailabilityByID]")

	availability, err := s.repository.GetAvailabilityByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return availability, ErrAvailabilityNotFound
		}
		return availability, err
	}

	return availability, nil
}

// UpdateAvailability updates an availability
func (s *Service) UpdateAvailability(id types.BinaryUUID, updateFn func(*models.Availability) error) error {
	s.logger.Info("[BookingService...UpdateAvailability]")

	// Get existing availability
	availability, err := s.repository.GetAvailabilityByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAvailabilityNotFound
		}
		return err
	}

	// Apply updates via callback function
	if err := updateFn(&availability); err != nil {
		return err
	}

	// Validate time range
	if availability.EndTime.Before(availability.StartTime) {
		return ErrInvalidTimeRange
	}

	// Save updated availability
	return s.repository.UpdateAvailability(&availability)
}

// DeleteAvailability deletes an availability
func (s *Service) DeleteAvailability(id types.BinaryUUID) error {
	s.logger.Info("[BookingService...DeleteAvailability]")

	// Check if availability exists
	_, err := s.repository.GetAvailabilityByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAvailabilityNotFound
		}
		return err
	}

	// Delete availability
	return s.repository.DeleteAvailability(id)
}

// ListAvailabilitiesByResourceID lists availabilities for a resource
func (s *Service) ListAvailabilitiesByResourceID(resourceID types.BinaryUUID) ([]models.Availability, error) {
	s.logger.Info("[BookingService...ListAvailabilitiesByResourceID]")

	// Check if resource exists
	_, err := s.repository.GetResourceByID(resourceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	return s.repository.ListAvailabilitiesByResourceID(resourceID)
}

// CheckResourceAvailability checks if a resource is available for a specific time period
func (s *Service) CheckResourceAvailability(resourceID types.BinaryUUID, start, end time.Time) (bool, error) {
	s.logger.Info("[BookingService...CheckResourceAvailability]")

	// Validate input parameters
	if end.Before(start) || start.Before(time.Now()) {
		return false, ErrInvalidTimeRange
	}

	// Check if resource exists
	_, err := s.repository.GetResourceByID(resourceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, ErrResourceNotFound
		}
		return false, err
	}

	// Check for overlapping bookings
	overlapping, err := s.repository.FindOverlappingBookings(resourceID, start, end)
	if err != nil {
		return false, err
	}

	if len(overlapping) > 0 {
		return false, nil
	}

	// Check if time falls within availability windows
	available, err := s.repository.IsAvailable(resourceID, start, end)
	if err != nil {
		return false, err
	}

	return available, nil
}

// -------------- Booking Service Methods --------------

// CreateBooking creates a new booking
func (s *Service) CreateBooking(booking *models.Booking) error {
	s.logger.Info("[BookingService...CreateBooking]")

	// Validate time range
	if booking.EndTime.Before(booking.StartTime) {
		return ErrInvalidTimeRange
	}

	// Check if booking is in the past
	if booking.StartTime.Before(time.Now()) {
		return ErrPastDateBooking
	}

	// Check availability first
	available, err := s.CheckResourceAvailability(booking.ResourceID, booking.StartTime, booking.EndTime)
	if err != nil {
		return err
	}

	if !available {
		return ErrResourceNotAvailable
	}

	// Generate UUID if not provided
	if booking.UUID.String() == (types.BinaryUUID{}).String() {
		id, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		booking.UUID = types.BinaryUUID(id)
	}

	// Set initial status if not provided
	if booking.Status == "" {
		booking.Status = "confirmed"
	}

	// Save to database
	return s.repository.CreateBooking(booking)
}

// GetBookingByID gets a booking by ID
func (s *Service) GetBookingByID(id types.BinaryUUID) (models.Booking, error) {
	s.logger.Info("[BookingService...GetBookingByID]")

	booking, err := s.repository.GetBookingByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return booking, ErrBookingNotFound
		}
		return booking, err
	}

	return booking, nil
}

// UpdateBooking updates a booking
func (s *Service) UpdateBooking(id types.BinaryUUID, updateFn func(*models.Booking) error) error {
	s.logger.Info("[BookingService...UpdateBooking]")

	// Get existing booking
	booking, err := s.repository.GetBookingByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBookingNotFound
		}
		return err
	}

	// Store original times to check availability if they change
	originalStart := booking.StartTime
	originalEnd := booking.EndTime

	// Apply updates via callback function
	if err := updateFn(&booking); err != nil {
		return err
	}

	// Validate time range
	if booking.EndTime.Before(booking.StartTime) {
		return ErrInvalidTimeRange
	}

	// Check status is valid
	if !isValidStatus(booking.Status) {
		return ErrInvalidBookingStatus
	}

	// If times changed, check availability
	if !booking.StartTime.Equal(originalStart) || !booking.EndTime.Equal(originalEnd) {
		// Check if booking is in the past
		if booking.StartTime.Before(time.Now()) {
			return ErrPastDateBooking
		}

		// For the availability check, we need to exclude the current booking
		// Get other overlapping bookings
		overlapping, err := s.repository.FindOverlappingBookings(booking.ResourceID, booking.StartTime, booking.EndTime)
		if err != nil {
			return err
		}

		// Filter out the current booking
		hasOverlap := false
		for _, b := range overlapping {
			if b.UUID != booking.UUID {
				hasOverlap = true
				break
			}
		}

		if hasOverlap {
			return ErrBookingOverlap
		}

		// Check if time falls within availability windows
		available, err := s.repository.IsAvailable(booking.ResourceID, booking.StartTime, booking.EndTime)
		if err != nil {
			return err
		}

		if !available {
			return ErrResourceNotAvailable
		}
	}

	// Save updated booking
	return s.repository.UpdateBooking(&booking)
}

// CancelBooking cancels a booking
func (s *Service) CancelBooking(id types.BinaryUUID) error {
	s.logger.Info("[BookingService...CancelBooking]")

	// Get existing booking
	booking, err := s.repository.GetBookingByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBookingNotFound
		}
		return err
	}

	// Set status to cancelled
	booking.Status = "cancelled"

	// Save updated booking
	return s.repository.UpdateBooking(&booking)
}

// ListBookings lists bookings with pagination and filtering
func (s *Service) ListBookings(page, limit int, filters map[string]interface{}) ([]models.Booking, int64, error) {
	s.logger.Info("[BookingService...ListBookings]")
	return s.repository.ListBookings(page, limit, filters)
}

// ListBookingsByUserID lists bookings for a specific user
func (s *Service) ListBookingsByUserID(userID types.BinaryUUID, page, limit int) ([]models.Booking, int64, error) {
	s.logger.Info("[BookingService...ListBookingsByUserID]")
	return s.repository.ListBookingsByUserID(userID, page, limit)
}

// Helper function to check if a booking status is valid
func isValidStatus(status string) bool {
	validStatuses := []string{"pending", "confirmed", "cancelled", "completed"}

	for _, s := range validStatuses {
		if status == s {
			return true
		}
	}

	return false
}
