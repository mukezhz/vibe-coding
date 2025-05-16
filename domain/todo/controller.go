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
		responses.HandleValidationError(ctx, c.logger, ErrTodoTitleRequired)
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
		responses.HandleValidationError(ctx, c.logger, ErrInvalidTodoID)
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
		responses.HandleValidationError(ctx, c.logger, ErrInvalidTodoID)
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
	if err != nil || limit < 1 {
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
