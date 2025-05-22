---
applyTo: "**/*.go"
---
# Error Handling Guide for Go Clean Architecture

This comprehensive guide explains how to handle errors consistently in our Go Clean Architecture project.

## Error Types and Hierarchy

Our project uses a structured error system defined in the `pkg/errorz` package:

### Base Error Types

- **APIError**: Base error type with status code, message, and error code
- **ValidationError**: For input validation failures
- **BusinessError**: For business rule violations
- **SystemError**: For internal system errors

### Common Predefined Errors

```go
// Common HTTP error types
var (
    ErrBadRequest   = NewAPIError(http.StatusBadRequest, "Bad Request", "BAD_REQUEST")
    ErrUnauthorized = NewAPIError(http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
    ErrForbidden    = NewAPIError(http.StatusForbidden, "Forbidden", "FORBIDDEN")
    ErrNotFound     = NewAPIError(http.StatusNotFound, "Not Found", "NOT_FOUND")
    ErrInternal     = NewAPIError(http.StatusInternalServerError, "Internal Server Error", "INTERNAL_ERROR")
)
```

## Domain-Specific Errors

Each domain should define its own error types in an `errorz.go` file:

```go
// domain/feature/errorz.go
package feature

import (
	"clean-architecture/pkg/errorz"
)

var (
	ErrInvalidFeatureID = errorz.ErrBadRequest.JoinError("Invalid Feature ID")
	ErrFeatureNotFound  = errorz.ErrNotFound.JoinError("Feature not found")
)

// Constructor functions for errors
func NewFeatureNotFoundError() *errorz.APIError {
	return ErrFeatureNotFound
}
```

## Response Handling

The `responses` package provides utility functions to standardize HTTP responses:

### Key Functions

1. **DetailResponse**: For single item responses
2. **ListResponse**: For collections without pagination
3. **ListResponseWithPagination**: For paginated collections
4. **SuccessResponse**: Simple success message
5. **ErrorResponse**: For error responses

### Error Handling Functions

1. **HandleValidationError**: For validation errors (400 Bad Request)
2. **HandleError**: General error handler that:
   - Handles custom `APIError` types
   - Handles `gorm.ErrRecordNotFound` with 404 Not Found
   - Logs and sends 500 Internal Server Error for unhandled errors
   - Captures unhandled exceptions using Sentry

## Examples

### Service Layer Error Handling

```go
// Service layer should map technical errors to domain errors
func (s *Service) GetByID(id types.BinaryUUID) (*models.Feature, error) {
    feature, err := s.repository.GetByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, NewFeatureNotFoundError()
        }
        return nil, err
    }
    return feature, nil
}
```

### Controller Layer Error Handling

```go
func (c *Controller) GetFeatureByID(ctx *gin.Context) {
    id, err := types.ParseUUID(ctx.Param("id"))
    if err != nil {
        responses.HandleError(c.logger, ctx, NewInvalidFeatureIDError())
        return
    }

    feature, err := c.service.GetByID(id)
    if err != nil {
        responses.HandleError(c.logger, ctx, err)
        return
    }

    // Success response handling...
}
```

## Handling Unhandled Exceptions with Sentry

To handle panics and unexpected errors, use Sentry integration:

```go
defer func() {
    if r := recover(); r != nil {
        err := errors.New("Unhandled exception occurred")
        utils.CurrentSentryService.CaptureException(err)
        c.JSON(500, gin.H{"error": "An unexpected error occurred"})
    }
}()
```

### Best Practices

1. Use `defer` and `recover` to handle panics in critical sections
2. Define specific error types for each domain
3. Map technical errors to domain-specific errors in the service layer
4. Use constructor functions for creating error instances
5. Use the `responses` package functions for consistent error responses
6. Include contextual information in error messages
