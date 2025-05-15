# Complete Guide to Creating REST APIs in Our Go Clean Architecture

This enhanced guide provides concrete examples, code snippets, and patterns to help LLMs and developers generate consistent REST API endpoints in our Go Clean Architecture project.

## Architecture Overview

Our project strictly follows a layered architecture with clear separation of concerns. Each layer has its own responsibility:

1. **Models Layer** (`domain/models/`) - Database entities with GORM annotations
2. **Repository Layer** (`domain/<feature>/repository.go`) - Database access logic
3. **Service Layer** (`domain/<feature>/service.go`) - Business logic
4. **Controller Layer** (`domain/<feature>/controller.go`) - Request handling and response formation
5. **Route Layer** (`domain/<feature>/route.go`) - Endpoint definitions
6. **Module Layer** (`domain/<feature>/module.go`) - Dependency injection via fx
7. **DTO Layer** (`domain/<feature>/dto.go`) - Data Transfer Objects for requests/responses
8. **Error Layer** (`domain/<feature>/errorz.go`) - Feature-specific custom errors

## Step-by-Step Guide with Examples

### Step 1: Define Data Models

First, create or update models in the `domain/models/` directory:

```go
// File: domain/models/todo.go
package models

import (
	"time"

	"clean-architecture/pkg/types"
)

// Todo represents the todo model in the database
type Todo struct {
	ID          types.BinaryUUID `json:"id" gorm:"type:binary(16);primary_key"`
	Title       string           `json:"title" gorm:"not null"`
	Description string           `json:"description"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}
```

### Step 2: Create DTOs (Data Transfer Objects)

Define request and response structures in a `dto.go` file:

```go
// File: domain/todo/dto.go
package todo

import (
	"clean-architecture/pkg/responses"
	"time"
)

// CreateTodoRequest DTO for creating a todo
type CreateTodoRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

// TodoResponse DTO for todo response
type TodoResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TodoListItem DTO for items in todo list
type TodoListItem struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// TodoListResponse DTO for paginated todo list
type TodoListResponse = responses.ListResponseType[TodoListItem]

// UpdateTodoRequest DTO for updating a todo
type UpdateTodoRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}
```

### Step 3: Define Custom Errors (Optional but Recommended)

Create feature-specific errors in an `errorz.go` file:

```go
// File: domain/todo/errorz.go
package todo

import (
	"clean-architecture/pkg/errorz"
)

// Error codes specific to the todo domain
const (
	ErrInvalidTodoID     = "INVALID_TODO_ID"
	ErrTodoNotFound      = "TODO_NOT_FOUND"
	ErrTodoTitleRequired = "TODO_TITLE_REQUIRED"
)

// NewInvalidTodoIDError returns a new error for invalid todo ID
func NewInvalidTodoIDError() error {
	return errorz.ErrBadRequest.JoinError("Invalid Todo ID")
}

// NewTodoNotFoundError returns a new error when todo is not found
func NewTodoNotFoundError() error {
	return errorz.ErrNotFound.JoinError("Todo not found")
}

// NewTodoTitleRequiredError returns a new error when todo title is missing
func NewTodoTitleRequiredError() error {
	return errorz.ErrBadRequest.JoinError("Todo title is required")
}
```

### Step 4: Implement Repository Layer

Create the repository that interfaces with the database:

```go
// File: domain/todo/repository.go
package todo

import (
	"clean-architecture/domain/models"
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/infrastructure"
	"clean-architecture/pkg/types"
)

// Repository database structure
type Repository struct {
	infrastructure.Database
	logger framework.Logger
}

// NewRepository creates a new todo repository
func NewRepository(db infrastructure.Database, logger framework.Logger) Repository {
	return Repository{db, logger}
}

// Create creates a new todo
func (r *Repository) Create(todo *models.Todo) error {
	r.logger.Info("[TodoRepository...Create]")
	return r.DB.Create(todo).Error
}

// GetByID gets a todo by ID
func (r *Repository) GetByID(todoID types.BinaryUUID) (todo models.Todo, err error) {
	r.logger.Info("[TodoRepository...GetByID]")
	return todo, r.DB.Where("id = ?", todoID).First(&todo).Error
}

// Update updates a todo
func (r *Repository) Update(todo *models.Todo) error {
	r.logger.Info("[TodoRepository...Update]")
	return r.DB.Save(todo).Error
}

// List returns todos with pagination
func (r *Repository) List(page, limit int) (todos []models.Todo, total int64, err error) {
	r.logger.Info("[TodoRepository...List]")

	offset := (page - 1) * limit

	// Get total count
	if err = r.DB.Model(&models.Todo{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get todos with pagination
	err = r.DB.Offset(offset).Limit(limit).Find(&todos).Error
	return todos, total, err
}
```

### Step 5: Implement Service Layer

Create the service that implements business logic:

```go
// File: domain/todo/service.go
package todo

import (
	"clean-architecture/domain/models"
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/types"
)

// Service service layer
type Service struct {
	logger     framework.Logger
	repository Repository
}

// NewService creates a new todo service
func NewService(
	logger framework.Logger,
	repository Repository,
) *Service {
	return &Service{
		logger:     logger,
		repository: repository,
	}
}

// Create creates a new todo
func (s Service) Create(todo *models.Todo) error {
	return s.repository.Create(todo)
}

// GetByID gets a todo by ID
func (s Service) GetByID(todoID types.BinaryUUID) (models.Todo, error) {
	todo, err := s.repository.GetByID(todoID)
	if err != nil {
		// Check if it's a "not found" error
		if err.Error() == "record not found" {
			return todo, NewTodoNotFoundError()
		}
		return todo, err
	}
	return todo, nil
}

// Update updates a todo
func (s Service) Update(todo *models.Todo) error {
	return s.repository.Update(todo)
}

// List returns todos with pagination
func (s Service) List(page, limit int) ([]models.Todo, int64, error) {
	return s.repository.List(page, limit)
}
```

### Step 6: Implement Controller Layer

Create the controller to handle HTTP requests:

```go
// File: domain/todo/controller.go
package todo

import (
	"net/http"
	"strconv"
	"time"

	"clean-architecture/domain/models"
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/responses"
	"clean-architecture/pkg/types"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Controller data type
type Controller struct {
	service *Service
	logger  framework.Logger
	env     *framework.Env
}

// NewController creates new todo controller
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

// CreateTodo creates a new todo
func (c *Controller) CreateTodo(ctx *gin.Context) {
	var req CreateTodoRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	// Validate title is not empty
	if req.Title == "" {
		responses.HandleValidationError(ctx, c.logger, NewTodoTitleRequiredError())
		return
	}

	// Create new todo with UUID
	todo := &models.Todo{
		ID:          types.BinaryUUID(uuid.New()),
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := c.service.Create(todo); err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	todoResponse := TodoResponse{
		ID:          todo.ID.String(),
		Title:       todo.Title,
		Description: todo.Description,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}

	responses.DetailResponse(
		ctx,
		http.StatusCreated,
		responses.DetailResponseType[TodoResponse]{
			Item:    todoResponse,
			Message: "success",
		},
	)
}

// GetTodoByID gets a single todo by ID
func (c *Controller) GetTodoByID(ctx *gin.Context) {
	todoID := ctx.Param("id")

	parsedID, err := types.ShouldParseUUID(todoID)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, NewInvalidTodoIDError())
		return
	}

	todo, err := c.service.GetByID(parsedID)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	todoResponse := TodoResponse{
		ID:          todo.ID.String(),
		Title:       todo.Title,
		Description: todo.Description,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}

	responses.DetailResponse(
		ctx,
		http.StatusOK,
		responses.DetailResponseType[TodoResponse]{
			Item:    todoResponse,
			Message: "success",
		},
	)
}

// UpdateTodo updates a todo by ID
func (c *Controller) UpdateTodo(ctx *gin.Context) {
	todoID := ctx.Param("id")

	parsedID, err := types.ShouldParseUUID(todoID)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, NewInvalidTodoIDError())
		return
	}

	// Get the existing todo first
	todo, err := c.service.GetByID(parsedID)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Parse update request
	var req UpdateTodoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	// Update fields if provided
	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	todo.UpdatedAt = time.Now()

	if err := c.service.Update(&todo); err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	todoResponse := TodoResponse{
		ID:          todo.ID.String(),
		Title:       todo.Title,
		Description: todo.Description,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}

	responses.DetailResponse(
		ctx,
		http.StatusOK,
		responses.DetailResponseType[TodoResponse]{
			Item:    todoResponse,
			Message: "success",
		},
	)
}

// FetchTodoWithPagination gets todos with pagination
func (c *Controller) FetchTodoWithPagination(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	todos, total, err := c.service.List(page, limit)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Convert to response format
	items := make([]TodoListItem, len(todos))
	for i, todo := range todos {
		items[i] = TodoListItem{
			ID:    todo.ID.String(),
			Title: todo.Title,
		}
	}

	// Check if there are more items
	hasNext := (int64(page*limit) < total)

	response := TodoListResponse{
		Items:   items,
		Message: "success",
		Pagination: responses.PaginationResponseType{
			Total:   total,
			HasNext: hasNext,
		},
	}

	responses.ListResponse(
		ctx,
		http.StatusOK,
		response,
	)
}
```

### Step 7: Set Up Routes

Define API endpoints in a route file:

```go
// File: domain/product/route.go
package product

import (
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/infrastructure"
)

// Route structure for products
type Route struct {
	logger     framework.Logger
	handler    infrastructure.Router
	controller *Controller
}

// NewRoute initializes product routes
// IMPORTANT: Return as POINTER
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

// RegisterRoute configures product API endpoints
func RegisterRoute(r *Route) {
	r.logger.Info("Setting up product routes")
	
	// Group all product routes under /api/products
	api := r.handler.Group("/api/products")
	
	// Define RESTful endpoints
	api.POST("", r.controller.CreateProduct)         // Create a product
	api.GET("", r.controller.ListProducts)           // List all products with pagination
	api.GET("/:id", r.controller.GetProductByID)     // Get product by ID
	api.PUT("/:id", r.controller.UpdateProduct)      // Update product by ID
	api.DELETE("/:id", r.controller.DeleteProduct)   // Delete product by ID
}
```

### Step 8: Configure Dependency Injection

Set up the module for dependency injection with fx:

```go
// File: domain/product/module.go
package product

import "go.uber.org/fx"

// Module provides product dependencies
var Module = fx.Module("product",
	fx.Options(
		fx.Provide(
			NewRepository,
			NewService,
			NewController,
			NewRoute,
		),
		fx.Invoke(RegisterRoute),
	),
)
```

### Step 9: Add Feature Module to Domain Module

Update `domain/module.go` to include your new feature:

```go
// File: domain/module.go
package domain

import (
	"clean-architecture/domain/organization"
	"clean-architecture/domain/product" // Add this import
	"clean-architecture/domain/todo"
	"clean-architecture/domain/user"
	"go.uber.org/fx"
)

// Module combines all domain modules
var Module = fx.Module("domain",
	fx.Options(
		organization.Module,
		todo.Module,
		user.Module,
		product.Module, // Add this module
	),
)
```

## Standards for API Responses

### Success Response Format

All successful responses should follow consistent formats depending on the type of data being returned.

Example for a single resource (DetailResponseType):

```json
{
  "item": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "Product Name",
    "price": 29.99,
    "description": "Product description",
    "sku": "PROD-123",
    "inventory": 100,
    "status": "active",
    "created_at": "2023-01-01T12:00:00Z",
    "updated_at": "2023-01-02T15:30:00Z"
  },
  "message": "success"
}
```

Example for a list of resources (ListResponseType):

```json
{
  "items": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "Product 1",
      "price": 29.99,
      "sku": "PROD-123",
      "inventory": 100,
      "status": "active"
    },
    {
      "id": "223e4567-e89b-12d3-a456-426614174001",
      "name": "Product 2",
      "price": 19.99,
      "sku": "PROD-456",
      "inventory": 50,
      "status": "active"
    }
  ],
  "message": "success",
  "pagination": {
    "total": 50,
    "has_next": true
  }
}
```

### Error Response Format

All error responses should follow this format:

```json
{
  "error": "Error message"
}
```

Example error responses:

```json
{
  "error": "Product not found"
}
```

For validation errors, the response might include more details:

```json
{
  "error": "Validation failed",
  "details": {
    "name": "The name field is required",
    "price": "The price must be greater than 0"
  }
}
```

## Common Data Types and GORM Tags

| Field Type         | Go Type          | GORM Tag Example                     |
| ------------------ | ---------------- | ------------------------------------ |
| Primary Key (UUID) | types.BinaryUUID | `gorm:"type:binary(16);primary_key"` |
| String (short)     | string           | `gorm:"size:255;not null"`           |
| String (long)      | string           | `gorm:"type:text"`                   |
| Integer            | int              | `gorm:"type:int;not null"`           |
| Decimal            | float64          | `gorm:"type:decimal(10,2);not null"` |
| Boolean            | bool             | `gorm:"default:false"`               |
| Foreign Key        | types.BinaryUUID | `gorm:"type:binary(16);index"`       |
| Date/Time          | time.Time        | `gorm:"type:datetime"`               |
| Enum               | string           | `gorm:"size:20;default:'active'"`    |
| Soft Delete        | gorm.DeletedAt   | `gorm:"index"`                       |

## Common Validation Tags for DTOs

| Validation Need | Binding Tag                         | Example Field                                                                     |
| --------------- | ----------------------------------- | --------------------------------------------------------------------------------- |
| Required Field  | `binding:"required"`                | `Name string \`json:"name" binding:"required"\``                                  |
| String Length   | `binding:"min=3,max=255"`           | `Title string \`json:"title" binding:"required,min=3,max=255"\``                  |
| Numeric Range   | `binding:"min=0,max=100"`           | `Quantity int \`json:"quantity" binding:"required,min=1"\``                       |
| Email           | `binding:"email"`                   | `Email string \`json:"email" binding:"required,email"\``                          |
| URL             | `binding:"url"`                     | `Website string \`json:"website" binding:"omitempty,url"\``                       |
| Enum Values     | `binding:"oneof=value1 value2"`     | `Status string \`json:"status" binding:"omitempty,oneof=active inactive draft"\`` |
| UUID Format     | `binding:"uuid4"`                   | `ID string \`json:"id" binding:"required,uuid4"\``                                |
| Optional Field  | `binding:"omitempty"`               | `Tags []string \`json:"tags" binding:"omitempty"\``                               |
| Conditional     | `binding:"required_if=Field value"` | `Code string \`json:"code" binding:"required_if=Type discount"\``                 |

## Testing REST APIs

### Recommended Testing Flow

1. **Unit Testing**:
   - Test repository methods in isolation
   - Test service methods with mocked repositories
   - Test controllers with mocked services

2. **Integration Testing**:
   - Test API endpoints with real database connections
   - Use transactions to rollback after each test

### Example Test Structure

```go
// File: domain/product/service_test.go
package product_test

import (
	"clean-architecture/domain/models"
	"clean-architecture/domain/product"
	"clean-architecture/pkg/framework"
	"testing"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of Repository
type MockRepository struct {
	mock.Mock
}

// Implement all Repository methods for the mock...

func TestCreateProduct(t *testing.T) {
	// Set up mock repository
	mockRepo := new(MockRepository)
	logger := framework.NewLogger()
	
	// Set up service with mock dependencies
	service := product.NewService(logger, mockRepo)
	
	// Define test cases
	testCases := []struct {
		name          string
		product       *models.Product
		setupMock     func()
		expectedError error
	}{
		{
			name: "Success case",
			product: &models.Product{
				Name:  "Test Product",
				Price: 19.99,
				SKU:   "TEST-SKU-001",
			},
			setupMock: func() {
				// Mock GetBySKU to return not found
				mockRepo.On("GetBySKU", "TEST-SKU-001").Return(models.Product{}, gorm.ErrRecordNotFound)
				// Mock Create to return nil error
				mockRepo.On("Create", mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		// Add more test cases...
	}
	
	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up mock
			tc.setupMock()
			
			// Call method being tested
			err := service.Create(tc.product)
			
			// Assert results
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			
			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}
```

## Common Pitfalls and Best Practices

1. **Return Types for Dependencies**
   - Repository: Return as VALUE (non-pointer)
   - Service, Controller, Route: Return as POINTERS

2. **Error Handling**
   - Use domain-specific errors defined in your `errorz.go` file
   - Use `responses.HandleError(logger, ctx, err)` for consistent error responses
   - Use `responses.HandleValidationError(logger, ctx, err)` for validation errors

3. **Validation**
   - Use `binding` tags on request DTOs
   - Use pointers for optional fields in update requests (allows distinguishing nil from zero values)

4. **UUID Handling**
   - Use `types.BinaryUUID` for database storage
   - Convert to/from string for API requests/responses

5. **Logging**
   - Log entry and exit points in each method
   - Use structured logging with meaningful context

6. **Transaction Management**
   - For operations affecting multiple entities, use transactions
   - Example in repository: `tx := r.DB.Begin(); defer func() { if r := recover(); r != nil || err != nil { tx.Rollback() } else { tx.Commit() } }()`

## Command Reference

| Make Command          | Description                           |
| --------------------- | ------------------------------------- |
| `make dev`            | Start the development server          |
| `make migrate-status` | Show migration status                 |
| `make migrate-diff`   | Generate migration from model changes |
| `make migrate-apply`  | Apply pending migrations              |
| `make test`           | Run tests                             |
| `make lint`           | Run linters                           |

## Example Prompt for LLM to Generate a New API

```
Please create a complete REST API for a [feature_name] resource with the following specs:

Model fields:
- ID (UUID)
- [field1] (string, required)
- [field2] (int, optional)
- [field3] (date, required)
- CreatedAt/UpdatedAt timestamps

Required CRUD endpoints:
- POST /api/[feature_name] - Create new [feature_name]
- GET /api/[feature_name]/:id - Get [feature_name] by ID
- GET /api/[feature_name] - List [feature_name]s with pagination
- PUT /api/[feature_name]/:id - Update [feature_name]
- DELETE /api/[feature_name]/:id - Delete [feature_name] (optional)

Please generate all required files following our clean architecture pattern based on the todo implementation:

1. Model in domain/models/[feature_name].go
   - Define the database entity with proper GORM tags
   - Use types.BinaryUUID for ID field
   - Include created_at and updated_at fields

2. DTO in domain/[feature_name]/dto.go
   - Create*Request struct with validation tags
   - Update*Request struct with pointer fields for partial updates
   - *Response struct for single item responses
   - *ListItem struct for simplified list responses
   - *ListResponse type using generic responses.ListResponseType

3. Custom errors in domain/[feature_name]/errorz.go
   - Define error constants
   - Create error constructor functions
   - Map to appropriate HTTP status codes

4. Repository in domain/[feature_name]/repository.go
   - Create, GetByID, Update, List methods
   - Proper logging with context
   - Efficient pagination for List operation

5. Service in domain/[feature_name]/service.go
   - Business logic layer
   - Input validation and error mapping
   - Return domain-specific errors

6. Controller in domain/[feature_name]/controller.go
   - Request binding and validation
   - Convert between DTO and model
   - Use responses.HandleError and responses.HandleValidationError
   - Implement consistent responses with DetailResponse and ListResponse

7. Route in domain/[feature_name]/route.go
   - Group routes under /api/[feature_name]
   - Map HTTP methods to controller functions

8. Module in domain/[feature_name]/module.go
   - Set up dependency injection with fx.Provide
   - Register routes with fx.Invoke

9. Update domain/module.go
   - Import the new feature module
   - Add it to the domain Module options

Please follow our established patterns from the todo implementation for validation, error handling, response formatting, and dependency injection.
```

By following this comprehensive guide, you should be able to create consistent REST APIs in our Go Clean Architecture project.
