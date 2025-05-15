---
mode: 'agent'
tools: []
description: 'Generate a complete booking system with availability check in Go Clean Architecture'
---

# Booking System Implementation Guide

You are an expert Go developer specialized in clean architecture. Your task is to implement a booking system with availability checking functionality following our established patterns.

## System Requirements

### Core Entities

1. **Resources**
   - Bookable items (e.g., rooms, equipment, services)
   - Each resource has attributes like name, description, capacity, location
   - Resources may have different availability schedules

2. **Availability**
   - Time slots when resources are available for booking
   - May include recurring schedules or custom availability periods
   - Blackout dates or maintenance periods

3. **Bookings**
   - Reservations of resources for specific time periods
   - Include user information, resource details, timestamps
   - Status tracking (confirmed, pending, cancelled)

4. **Users**
   - Information about who made the booking
   - May include roles (admin, regular user) for authorization

### Required Functionality

1. **Resource Management**
   - CRUD operations for resources
   - Resource categorization and filtering

2. **Availability Checking**
   - Check if a resource is available for a specific time slot
   - Handle timezone considerations
   - Prevent double-booking

3. **Booking Process**
   - Create bookings with validation
   - Update booking status
   - Cancel or reschedule bookings
   - Send notifications upon booking changes

4. **Reporting**
   - View upcoming bookings
   - Generate utilization reports
   - Filter bookings by date range, user, resource

## Implementation Guide

Follow our clean architecture pattern to implement the booking system:

0. **Bruno Docs**
   - Use the Bruno Docs for reference on our clean architecture patterns and practices.
   - Add bruno docs for api for testing.

1. **Domain Models**
   - Create models in `domain/models/` with appropriate relationships
   - Example models: Resource, Availability, Booking, User

2. **Database Schema**
   - Design efficient schema with proper indexes
   - Consider using JSON fields for flexible attributes

3. **API Endpoints**
   - Follow RESTful conventions
   - Implement proper validation and error handling
   - Use consistent response formats with our generic response types

4. **Business Logic**
   - Implement robust availability checking algorithm
   - Consider concurrent booking scenarios and race conditions
   - Add business rules for booking limitations, lead times, etc.

5. **Testing**
   - Write comprehensive unit tests
   - Test edge cases like overlapping bookings
   - Add integration tests for booking workflows

## Example API Endpoints

### Resource Management
- `GET /api/resources` - List all resources with filtering
- `GET /api/resources/:id` - Get resource details
- `POST /api/resources` - Create a new resource
- `PUT /api/resources/:id` - Update a resource
- `DELETE /api/resources/:id` - Delete a resource

### Availability
- `GET /api/resources/:id/availability?start=<timestamp>&end=<timestamp>` - Check resource availability
- `POST /api/resources/:id/availability` - Set availability for a resource
- `GET /api/availability?resource_ids=1,2,3&start=<timestamp>&end=<timestamp>` - Check multiple resources

### Bookings
- `GET /api/bookings` - List bookings with filtering
- `GET /api/bookings/:id` - Get booking details
- `POST /api/bookings` - Create a new booking
- `PUT /api/bookings/:id` - Update a booking
- `DELETE /api/bookings/:id` - Cancel a booking
- `GET /api/users/:id/bookings` - Get user's bookings

## Implementation Details

### Models Example

```go
// Resource model
type Resource struct {
    gorm.Model
    UUID        types.BinaryUUID `json:"uuid" gorm:"index;notnull;unique"`
    Name        string           `json:"name" gorm:"size:255;not null"`
    Description string           `json:"description" gorm:"type:text"`
    Type        string           `json:"type" gorm:"size:50;not null"`
    Capacity    int              `json:"capacity" gorm:"default:1"`
    Location    string           `json:"location" gorm:"size:255"`
    Attributes  datatypes.JSON   `json:"attributes" gorm:"type:json"`
}

// Availability model
type Availability struct {
    gorm.Model
    UUID        types.BinaryUUID `json:"uuid" gorm:"index;notnull;unique"`
    ResourceID  types.BinaryUUID `json:"resource_id" gorm:"index;not null"`
    StartTime   time.Time        `json:"start_time" gorm:"not null;index"`
    EndTime     time.Time        `json:"end_time" gorm:"not null;index"`
    IsRecurring bool             `json:"is_recurring" gorm:"default:false"`
    RecurRule   string           `json:"recur_rule" gorm:"size:255"`
}

// Booking model
type Booking struct {
    gorm.Model
    UUID        types.BinaryUUID `json:"uuid" gorm:"index;notnull;unique"`
    ResourceID  types.BinaryUUID `json:"resource_id" gorm:"index;not null"`
    UserID      types.BinaryUUID `json:"user_id" gorm:"index;not null"`
    StartTime   time.Time        `json:"start_time" gorm:"not null;index"`
    EndTime     time.Time        `json:"end_time" gorm:"not null;index"`
    Status      string           `json:"status" gorm:"size:50;default:'pending'"`
    Notes       string           `json:"notes" gorm:"type:text"`
    Reference   string           `json:"reference" gorm:"size:100"`
}
```

### Availability Checking Logic

```go
// Service layer availability check
func (s *Service) CheckResourceAvailability(resourceID types.BinaryUUID, start, end time.Time) (bool, error) {
    // Validate input parameters
    if end.Before(start) || start.Before(time.Now()) {
        return false, errors.New("invalid time range")
    }
    
    // Check if resource exists
    _, err := s.resourceRepo.GetByID(resourceID)
    if err != nil {
        return false, err
    }
    
    // Check for overlapping bookings
    overlapping, err := s.bookingRepo.FindOverlappingBookings(resourceID, start, end)
    if err != nil {
        return false, err
    }
    
    if len(overlapping) > 0 {
        return false, nil
    }
    
    // Check if time falls within availability windows
    available, err := s.availabilityRepo.IsAvailable(resourceID, start, end)
    if err != nil {
        return false, err
    }
    
    return available, nil
}
```

### Booking Creation Logic

```go
// Service layer booking creation
func (s *Service) CreateBooking(booking *models.Booking) error {
    // Check availability first
    available, err := s.CheckResourceAvailability(booking.ResourceID, booking.StartTime, booking.EndTime)
    if err != nil {
        return err
    }
    
    if !available {
        return NewResourceNotAvailableError()
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
    if err := s.bookingRepo.Create(booking); err != nil {
        return err
    }
    
    // Send notification
    // ...
    
    return nil
}
```

## Best Practices

1. **Concurrency Control**
   - Use database transactions for booking creation
   - Consider implementing optimistic locking
   - Handle race conditions for concurrent booking attempts

2. **Performance Optimization**
   - Index time-based queries
   - Cache availability for popular resources
   - Implement pagination for booking lists

3. **Business Rules**
   - Enforce minimum and maximum booking duration
   - Implement advanced notice requirements
   - Handle cancellation policies and fees

4. **Error Handling**
   - Create domain-specific errors (ResourceNotAvailableError, BookingOverlapError)
   - Return clear error messages for validation failures
   - Log detailed error information for debugging

Follow our project's established patterns for controller implementation, response formatting, and error handling as documented in the AddingEndpoints guide.
