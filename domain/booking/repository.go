package booking

import (
	"time"

	"clean-architecture/domain/models"
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/infrastructure"
	"clean-architecture/pkg/types"
)

// Repository handles database operations for resources, availability, and bookings
type Repository struct {
	infrastructure.Database
	logger framework.Logger
}

// NewRepository creates a new booking repository
func NewRepository(db infrastructure.Database, logger framework.Logger) Repository {
	return Repository{db, logger}
}

// -------------- Resource Repository Methods --------------

// CreateResource adds a new resource to the database
func (r Repository) CreateResource(resource *models.Resource) error {
	r.logger.Info("[BookingRepository...CreateResource]")
	return r.DB.Create(resource).Error
}

// GetResourceByID retrieves a resource by ID
func (r Repository) GetResourceByID(id types.BinaryUUID) (models.Resource, error) {
	r.logger.Info("[BookingRepository...GetResourceByID]")
	var resource models.Resource
	err := r.DB.Where("uuid = ?", id).First(&resource).Error
	return resource, err
}

// UpdateResource updates a resource
func (r Repository) UpdateResource(resource *models.Resource) error {
	r.logger.Info("[BookingRepository...UpdateResource]")
	return r.DB.Save(resource).Error
}

// DeleteResource deletes a resource
func (r Repository) DeleteResource(id types.BinaryUUID) error {
	r.logger.Info("[BookingRepository...DeleteResource]")
	return r.DB.Where("uuid = ?", id).Delete(&models.Resource{}).Error
}

// ListResources returns resources with pagination and filtering
func (r Repository) ListResources(page, limit int, filters map[string]interface{}) ([]models.Resource, int64, error) {
	r.logger.Info("[BookingRepository...ListResources]")
	var resources []models.Resource
	var total int64

	query := r.DB

	// Apply filters if any
	for key, value := range filters {
		if value != nil && value != "" {
			query = query.Where(key+" = ?", value)
		}
	}

	// Get total count
	if err := query.Model(&models.Resource{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&resources).Error

	return resources, total, err
}

// -------------- Availability Repository Methods --------------

// CreateAvailability adds a new availability to the database
func (r Repository) CreateAvailability(availability *models.Availability) error {
	r.logger.Info("[BookingRepository...CreateAvailability]")
	return r.DB.Create(availability).Error
}

// GetAvailabilityByID retrieves an availability by ID
func (r Repository) GetAvailabilityByID(id types.BinaryUUID) (models.Availability, error) {
	r.logger.Info("[BookingRepository...GetAvailabilityByID]")
	var availability models.Availability
	err := r.DB.Where("uuid = ?", id).First(&availability).Error
	return availability, err
}

// UpdateAvailability updates an availability
func (r Repository) UpdateAvailability(availability *models.Availability) error {
	r.logger.Info("[BookingRepository...UpdateAvailability]")
	return r.DB.Save(availability).Error
}

// DeleteAvailability deletes an availability
func (r Repository) DeleteAvailability(id types.BinaryUUID) error {
	r.logger.Info("[BookingRepository...DeleteAvailability]")
	return r.DB.Where("uuid = ?", id).Delete(&models.Availability{}).Error
}

// ListAvailabilitiesByResourceID returns availabilities for a resource
func (r Repository) ListAvailabilitiesByResourceID(resourceID types.BinaryUUID) ([]models.Availability, error) {
	r.logger.Info("[BookingRepository...ListAvailabilitiesByResourceID]")
	var availabilities []models.Availability
	err := r.DB.Where("resource_id = ?", resourceID).Find(&availabilities).Error
	return availabilities, err
}

// IsAvailable checks if a resource is available for a specific time period
func (r Repository) IsAvailable(resourceID types.BinaryUUID, start, end time.Time) (bool, error) {
	r.logger.Info("[BookingRepository...IsAvailable]")

	// Check for any availability windows that cover the requested time
	var count int64
	err := r.DB.Model(&models.Availability{}).
		Where("resource_id = ? AND start_time <= ? AND end_time >= ?", resourceID, start, end).
		Count(&count).Error

	return count > 0, err
}

// -------------- Booking Repository Methods --------------

// CreateBooking adds a new booking to the database
func (r Repository) CreateBooking(booking *models.Booking) error {
	r.logger.Info("[BookingRepository...CreateBooking]")
	return r.DB.Create(booking).Error
}

// GetBookingByID retrieves a booking by ID
func (r Repository) GetBookingByID(id types.BinaryUUID) (models.Booking, error) {
	r.logger.Info("[BookingRepository...GetBookingByID]")
	var booking models.Booking
	err := r.DB.Where("uuid = ?", id).First(&booking).Error
	return booking, err
}

// UpdateBooking updates a booking
func (r Repository) UpdateBooking(booking *models.Booking) error {
	r.logger.Info("[BookingRepository...UpdateBooking]")
	return r.DB.Save(booking).Error
}

// DeleteBooking cancels a booking
func (r Repository) DeleteBooking(id types.BinaryUUID) error {
	r.logger.Info("[BookingRepository...DeleteBooking]")
	// Soft delete for bookings
	return r.DB.Model(&models.Booking{}).Where("uuid = ?", id).Update("status", "cancelled").Error
}

// ListBookings returns bookings with pagination and filtering
func (r Repository) ListBookings(page, limit int, filters map[string]interface{}) ([]models.Booking, int64, error) {
	r.logger.Info("[BookingRepository...ListBookings]")
	var bookings []models.Booking
	var total int64

	query := r.DB

	// Apply filters if any
	for key, value := range filters {
		if value != nil && value != "" {
			query = query.Where(key+" = ?", value)
		}
	}

	// Get total count
	if err := query.Model(&models.Booking{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("start_time ASC").Find(&bookings).Error

	return bookings, total, err
}

// FindOverlappingBookings finds bookings that overlap with a time range for a resource
func (r Repository) FindOverlappingBookings(resourceID types.BinaryUUID, start, end time.Time) ([]models.Booking, error) {
	r.logger.Info("[BookingRepository...FindOverlappingBookings]")
	var bookings []models.Booking

	// Time range overlap query
	// (StartA <= EndB) AND (EndA >= StartB)
	err := r.DB.Where("resource_id = ? AND start_time <= ? AND end_time >= ? AND status != 'cancelled'",
		resourceID, end, start).Find(&bookings).Error

	return bookings, err
}

// ListBookingsByUserID returns bookings for a specific user
func (r Repository) ListBookingsByUserID(userID types.BinaryUUID, page, limit int) ([]models.Booking, int64, error) {
	r.logger.Info("[BookingRepository...ListBookingsByUserID]")
	var bookings []models.Booking
	var total int64

	// Get total count
	if err := r.DB.Model(&models.Booking{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	err := r.DB.Where("user_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Order("start_time ASC").
		Find(&bookings).Error

	return bookings, total, err
}
