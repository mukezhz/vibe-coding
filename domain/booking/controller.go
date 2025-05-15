package booking

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"clean-architecture/domain/models"
	"clean-architecture/pkg/errorz"
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/responses"
	"clean-architecture/pkg/types"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Controller handles HTTP requests for the booking system
type Controller struct {
	service *Service
	logger  framework.Logger
	env     *framework.Env
}

// NewController creates a new booking controller
func NewController(
	service *Service,
	logger framework.Logger,
	env *framework.Env,
) *Controller {
	return &Controller{
		service: service,
		logger:  logger,
		env:     env,
	}
}

// -------------- Resource Controllers --------------

// CreateResource handles the create resource request
func (c *Controller) CreateResource(ctx *gin.Context) {
	c.logger.Info("[BookingController...CreateResource]")

	var req ResourceCreateDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	// Convert attributes to JSON
	var attributes datatypes.JSON
	if req.Attributes != nil {
		var err error
		attributesBytes, err := json.Marshal(req.Attributes)
		if err != nil {
			responses.HandleError(ctx, c.logger, err)
			return
		}
		err = attributes.UnmarshalJSON(attributesBytes)
		if err != nil {
			responses.HandleError(ctx, c.logger, err)
			return
		}
	}

	// Convert request to model
	resource := models.Resource{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Capacity:    req.Capacity,
		Location:    req.Location,
		Attributes:  attributes,
	}

	// Create resource
	if err := c.service.CreateResource(&resource); err != nil {
		c.logger.Errorf("[BookingController...CreateResource] Error: %v", err)
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Convert to response DTO
	response := ResourceToDTO(&resource)

	responses.DetailResponse(
		ctx,
		http.StatusCreated,
		responses.DetailResponseType[ResourceResponseDTO]{
			Item:    response,
			Message: "Resource created successfully",
		},
	)
}

// GetResourceByID handles the get resource by ID request
func (c *Controller) GetResourceByID(ctx *gin.Context) {
	c.logger.Info("[BookingController...GetResourceByID]")

	// Parse ID parameter
	idParam := ctx.Param("id")
	parsedID, err := types.ShouldParseUUID(idParam)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, errorz.ErrBadRequest)
		return
	}

	// Get resource
	resource, err := c.service.GetResourceByID(parsedID)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Convert to response DTO
	response := ResourceToDTO(&resource)

	responses.DetailResponse(
		ctx,
		http.StatusOK,
		responses.DetailResponseType[ResourceResponseDTO]{
			Item:    response,
			Message: "success",
		},
	)
}

// UpdateResource handles the update resource request
func (c *Controller) UpdateResource(ctx *gin.Context) {
	c.logger.Info("[BookingController...UpdateResource]")

	// Parse ID parameter
	idParam := ctx.Param("id")
	parsedID, err := types.ShouldParseUUID(idParam)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, errorz.ErrBadRequest)
		return
	}

	// Parse request body
	var req ResourceUpdateDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	// Update resource
	err = c.service.UpdateResource(parsedID, func(resource *models.Resource) error {
		if req.Name != "" {
			resource.Name = req.Name
		}
		if req.Description != "" {
			resource.Description = req.Description
		}
		if req.Type != "" {
			resource.Type = req.Type
		}
		if req.Capacity != 0 {
			resource.Capacity = req.Capacity
		}
		if req.Location != "" {
			resource.Location = req.Location
		}
		if req.Attributes != nil {
			attributesBytes, err := json.Marshal(req.Attributes)
			if err != nil {
				return err
			}
			var jsonData datatypes.JSON
			if err := jsonData.UnmarshalJSON(attributesBytes); err != nil {
				return err
			}
			resource.Attributes = jsonData
		}

		return nil
	})

	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Get updated resource
	resource, err := c.service.GetResourceByID(parsedID)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Convert to response DTO
	response := ResourceToDTO(&resource)

	responses.DetailResponse(
		ctx,
		http.StatusOK,
		responses.DetailResponseType[ResourceResponseDTO]{
			Item:    response,
			Message: "Resource updated successfully",
		},
	)
}

// DeleteResource handles the delete resource request
func (c *Controller) DeleteResource(ctx *gin.Context) {
	c.logger.Info("[BookingController...DeleteResource]")

	// Parse ID parameter
	idParam := ctx.Param("id")
	parsedID, err := types.ShouldParseUUID(idParam)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, errorz.ErrBadRequest)
		return
	}

	// Delete resource
	if err := c.service.DeleteResource(parsedID); err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Return success response (no content)
	ctx.Status(http.StatusNoContent)
}

// ListResources handles the list resources request with pagination
func (c *Controller) ListResources(ctx *gin.Context) {
	c.logger.Info("[BookingController...ListResources]")

	// Parse pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Parse filters
	filters := make(map[string]interface{})
	if resourceType := ctx.Query("type"); resourceType != "" {
		filters["type"] = resourceType
	}
	if location := ctx.Query("location"); location != "" {
		filters["location"] = location
	}
	if capacity := ctx.Query("capacity"); capacity != "" {
		capInt, err := strconv.Atoi(capacity)
		if err == nil {
			filters["capacity"] = capInt
		}
	}

	// Get resources
	resources, total, err := c.service.ListResources(page, limit, filters)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Convert to response format
	items := make([]ResourceResponseDTO, len(resources))
	for i, resource := range resources {
		items[i] = ResourceToDTO(&resource)
	}

	// Create paginated response
	responses.ListResponse(
		ctx,
		http.StatusOK,
		responses.ListResponseType[ResourceResponseDTO]{
			Items: items,
			Pagination: responses.PaginationResponseType{
				Total:   total,
				HasNext: int64(page*limit) < total,
			},
			Message: "Resources retrieved successfully",
		},
	)
}

// -------------- Availability Controllers --------------

// CreateAvailability handles the create availability request
func (c *Controller) CreateAvailability(ctx *gin.Context) {
	c.logger.Info("[BookingController...CreateAvailability]")

	// Parse resource ID parameter
	resourceIDParam := ctx.Param("id")
	resourceID, err := uuid.Parse(resourceIDParam)
	if err != nil {
		responses.HandleError(ctx, c.logger, errorz.ErrBadRequest)
		return
	}

	var req AvailabilityCreateDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	// Convert request to model
	availability := models.Availability{
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		IsRecurring: req.IsRecurring,
		RecurRule:   req.RecurRule,
	}

	// Create availability
	if err := c.service.CreateAvailability(types.BinaryUUID(resourceID), &availability); err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Convert to response DTO
	response := AvailabilityToDTO(&availability)

	responses.DetailResponse(
		ctx,
		http.StatusCreated,
		responses.DetailResponseType[AvailabilityResponseDTO]{
			Item:    response,
			Message: "Availability created successfully",
		},
	)
}

// CheckResourceAvailability handles the check resource availability request
func (c *Controller) CheckResourceAvailability(ctx *gin.Context) {
	c.logger.Info("[BookingController...CheckResourceAvailability]")

	// Parse resource ID parameter
	resourceIDParam := ctx.Param("id")
	resourceID, err := uuid.Parse(resourceIDParam)
	if err != nil {
		responses.HandleError(ctx, c.logger, errorz.ErrBadRequest)
		return
	}

	// Parse query parameters
	var query AvailabilityCheckDTO
	if err := ctx.ShouldBindQuery(&query); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	// Check resource availability
	available, err := c.service.CheckResourceAvailability(
		types.BinaryUUID(resourceID),
		query.StartTime,
		query.EndTime,
	)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Return availability response
	response := AvailabilityCheckResponseDTO{
		Available: available,
	}

	responses.DetailResponse(
		ctx,
		http.StatusOK,
		responses.DetailResponseType[AvailabilityCheckResponseDTO]{
			Item:    response,
			Message: "Availability check completed",
		},
	)
}

// ListResourceAvailabilities handles listing availabilities for a resource
func (c *Controller) ListResourceAvailabilities(ctx *gin.Context) {
	c.logger.Info("[BookingController...ListResourceAvailabilities]")

	// Parse resource ID parameter
	resourceIDParam := ctx.Param("id")
	resourceID, err := uuid.Parse(resourceIDParam)
	if err != nil {
		responses.HandleError(ctx, c.logger, errorz.ErrBadRequest)
		return
	}

	// Get availabilities
	availabilities, err := c.service.ListAvailabilitiesByResourceID(types.BinaryUUID(resourceID))
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Convert to response format
	items := make([]AvailabilityResponseDTO, len(availabilities))
	for i, availability := range availabilities {
		items[i] = AvailabilityToDTO(&availability)
	}

	responses.ListResponse(
		ctx,
		http.StatusOK,
		responses.ListResponseType[AvailabilityResponseDTO]{
			Items:   items,
			Message: "Availabilities retrieved successfully",
			Pagination: responses.PaginationResponseType{
				Total:   int64(len(items)),
				HasNext: false,
			},
		},
	)
}

// CheckMultipleResourcesAvailability handles checking availability for multiple resources
func (c *Controller) CheckMultipleResourcesAvailability(ctx *gin.Context) {
	c.logger.Info("[BookingController...CheckMultipleResourcesAvailability]")

	// Parse resource IDs
	resourceIDsParam := ctx.QueryArray("resource_ids")
	if len(resourceIDsParam) == 0 {
		responses.HandleError(ctx, c.logger, errorz.ErrBadRequest)
		return
	}

	// Parse query parameters
	startStr := ctx.Query("start")
	endStr := ctx.Query("end")

	if startStr == "" || endStr == "" {
		responses.HandleError(ctx, c.logger, errorz.ErrBadRequest)
		return
	}

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		responses.HandleError(ctx, c.logger, errorz.ErrBadRequest)
		return
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		responses.HandleError(ctx, c.logger, errorz.ErrBadRequest)
		return
	}

	// Check availability for each resource
	results := make(map[string]bool)

	for _, idStr := range resourceIDsParam {
		id, err := uuid.Parse(idStr)
		if err != nil {
			// Skip invalid IDs
			continue
		}

		available, err := c.service.CheckResourceAvailability(types.BinaryUUID(id), start, end)
		if err != nil {
			// Skip resources with errors
			continue
		}

		results[idStr] = available
	}

	// Return results
	response := struct {
		Results map[string]bool `json:"results"`
	}{
		Results: results,
	}

	responses.DetailResponse(
		ctx,
		http.StatusOK,
		responses.DetailResponseType[interface{}]{
			Item:    response,
			Message: "Availability check completed",
		},
	)
}

// -------------- Booking Controllers --------------

// CreateBooking handles the create booking request
func (c *Controller) CreateBooking(ctx *gin.Context) {
	c.logger.Info("[BookingController...CreateBooking]")

	var req BookingCreateDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	// Get user ID from context
	userIDStr := ctx.GetString("user_id")
	if userIDStr == "" {
		responses.HandleError(ctx, c.logger, errorz.ErrUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		responses.HandleError(ctx, c.logger, errorz.ErrUnauthorized)
		return
	}

	// Convert request to model
	booking := models.Booking{
		ResourceID: req.ResourceID,
		UserID:     types.BinaryUUID(userID),
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		Notes:      req.Notes,
		Reference:  req.Reference,
	}

	// Create booking
	if err := c.service.CreateBooking(&booking); err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Convert to response DTO
	response := BookingToDTO(&booking)

	responses.DetailResponse(
		ctx,
		http.StatusCreated,
		responses.DetailResponseType[BookingResponseDTO]{
			Item:    response,
			Message: "Booking created successfully",
		},
	)
}

// GetBookingByID handles the get booking by ID request
func (c *Controller) GetBookingByID(ctx *gin.Context) {
	c.logger.Info("[BookingController...GetBookingByID]")

	// Parse ID parameter
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		responses.HandleError(ctx, c.logger, errorz.ErrBadRequest)
		return
	}

	// Get booking
	booking, err := c.service.GetBookingByID(types.BinaryUUID(id))
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Authorization check: user can only see their own bookings unless they're an admin
	// TODO: Implement proper auth check with roles
	userIDStr := ctx.GetString("user_id")
	if userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err == nil && booking.UserID != types.BinaryUUID(userID) {
			// Check if user has admin role
			isAdmin := ctx.GetBool("is_admin") // Assuming this is set by auth middleware
			if !isAdmin {
				responses.HandleError(ctx, c.logger, errorz.ErrForbidden)
				return
			}
		}
	}

	// Convert to response DTO
	response := BookingToDTO(&booking)

	responses.DetailResponse(
		ctx,
		http.StatusOK,
		responses.DetailResponseType[BookingResponseDTO]{
			Item:    response,
			Message: "Booking retrieved successfully",
		},
	)
}

// UpdateBooking handles the update booking request
func (c *Controller) UpdateBooking(ctx *gin.Context) {
	c.logger.Info("[BookingController...UpdateBooking]")

	// Parse ID parameter
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		responses.HandleError(ctx, c.logger, errorz.ErrBadRequest)
		return
	}

	// Get booking to check authorization
	booking, err := c.service.GetBookingByID(types.BinaryUUID(id))
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Authorization check: user can only update their own bookings unless they're an admin
	userIDStr := ctx.GetString("user_id")
	if userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err == nil && booking.UserID != types.BinaryUUID(userID) {
			// Check if user has admin role
			isAdmin := ctx.GetBool("is_admin") // Assuming this is set by auth middleware
			if !isAdmin {
				responses.HandleError(ctx, c.logger, errorz.ErrForbidden)
				return
			}
		}
	}

	// Parse request body
	var req BookingUpdateDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	// Update booking
	err = c.service.UpdateBooking(types.BinaryUUID(id), func(booking *models.Booking) error {
		// Only update fields that were provided
		timeChanged := false

		if !req.StartTime.IsZero() {
			booking.StartTime = req.StartTime
			timeChanged = true
		}

		if !req.EndTime.IsZero() {
			booking.EndTime = req.EndTime
			timeChanged = true
		}

		// Only allow status updates if times didn't change
		if req.Status != "" && !timeChanged {
			booking.Status = req.Status
		}

		if req.Notes != "" {
			booking.Notes = req.Notes
		}

		if req.Reference != "" {
			booking.Reference = req.Reference
		}

		return nil
	})

	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Get updated booking
	updatedBooking, err := c.service.GetBookingByID(types.BinaryUUID(id))
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Convert to response DTO
	response := BookingToDTO(&updatedBooking)

	responses.DetailResponse(
		ctx,
		http.StatusOK,
		responses.DetailResponseType[BookingResponseDTO]{
			Item:    response,
			Message: "Booking updated successfully",
		},
	)
}

// CancelBooking handles the cancel booking request
func (c *Controller) CancelBooking(ctx *gin.Context) {
	c.logger.Info("[BookingController...CancelBooking]")

	// Parse ID parameter
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		responses.HandleError(ctx, c.logger, errorz.ErrBadRequest)
		return
	}

	// Get booking to check authorization
	booking, err := c.service.GetBookingByID(types.BinaryUUID(id))
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Authorization check: user can only cancel their own bookings unless they're an admin
	userIDStr := ctx.GetString("user_id")
	if userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err == nil && booking.UserID != types.BinaryUUID(userID) {
			// Check if user has admin role
			isAdmin := ctx.GetBool("is_admin") // Assuming this is set by auth middleware
			if !isAdmin {
				responses.HandleError(ctx, c.logger, errorz.ErrForbidden)
				return
			}
		}
	}

	// Cancel booking
	if err := c.service.CancelBooking(types.BinaryUUID(id)); err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Return success response (no content)
	ctx.Status(http.StatusNoContent)
}

// ListBookings handles listing bookings with filtering
func (c *Controller) ListBookings(ctx *gin.Context) {
	c.logger.Info("[BookingController...ListBookings]")

	// Parse pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Parse filters
	filters := make(map[string]interface{})

	// If user is not admin, restrict to their own bookings
	isAdmin := ctx.GetBool("is_admin") // Assuming this is set by auth middleware
	if !isAdmin {
		userIDStr := ctx.GetString("user_id")
		if userIDStr == "" {
			responses.HandleError(ctx, c.logger, errorz.ErrUnauthorized)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			responses.HandleError(ctx, c.logger, errorz.ErrBadRequest)
			return
		}

		filters["user_id"] = types.BinaryUUID(userID)
	} else {
		// Allow filtering by resource and user for admins
		if resourceIDStr := ctx.Query("resource_id"); resourceIDStr != "" {
			resourceID, err := uuid.Parse(resourceIDStr)
			if err == nil {
				filters["resource_id"] = types.BinaryUUID(resourceID)
			}
		}

		if userIDStr := ctx.Query("user_id"); userIDStr != "" {
			userID, err := uuid.Parse(userIDStr)
			if err == nil {
				filters["user_id"] = types.BinaryUUID(userID)
			}
		}
	}

	// Add common filters
	if status := ctx.Query("status"); status != "" {
		filters["status"] = status
	}

	// Get bookings
	bookings, total, err := c.service.ListBookings(page, limit, filters)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Convert to response format
	items := make([]BookingResponseDTO, len(bookings))
	for i, booking := range bookings {
		items[i] = BookingToDTO(&booking)
	}

	// Create paginated response
	responses.ListResponse(
		ctx,
		http.StatusOK,
		responses.ListResponseType[BookingResponseDTO]{
			Items:   items,
			Message: "Bookings retrieved successfully",
			Pagination: responses.PaginationResponseType{
				Total:   total,
				HasNext: int64(page*limit) < total,
			},
		},
	)
}

// ListUserBookings handles listing bookings for a specific user
func (c *Controller) ListUserBookings(ctx *gin.Context) {
	c.logger.Info("[BookingController...ListUserBookings]")

	// Parse user ID parameter
	userIDParam := ctx.Param("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		responses.HandleError(ctx, c.logger, errorz.ErrBadRequest)
		return
	}

	// Authorization check: user can only see their own bookings unless they're an admin
	requestingUserID := ctx.GetString("user_id")
	isAdmin := ctx.GetBool("is_admin") // Assuming this is set by auth middleware

	if requestingUserID != userIDParam && !isAdmin {
		responses.HandleError(ctx, c.logger, errorz.ErrForbidden)
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Get bookings
	bookings, total, err := c.service.ListBookingsByUserID(types.BinaryUUID(userID), page, limit)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Convert to response format
	items := make([]BookingResponseDTO, len(bookings))
	for i, booking := range bookings {
		items[i] = BookingToDTO(&booking)
	}

	// Calculate if there are more pages
	hasNext := int64(page*limit) < total

	// Create paginated response
	responses.ListResponse(
		ctx,
		http.StatusOK,
		responses.ListResponseType[BookingResponseDTO]{
			Items:   items,
			Message: "Bookings retrieved successfully",
			Pagination: responses.PaginationResponseType{
				Total:   total,
				HasNext: hasNext,
			},
		},
	)
}
