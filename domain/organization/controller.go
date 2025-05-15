package organization

import (
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/responses"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Controller handles HTTP requests for organizations
type Controller struct {
	service *Service
	logger  framework.Logger
}

// NewController creates a new organization controller
func NewController(
	service *Service,
	logger framework.Logger,
) *Controller {
	return &Controller{service, logger}
}

// Create handles the creation of a new organization
func (c *Controller) Create(ctx *gin.Context) {
	var request CreateOrganizationRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		responses.HandleError(ctx, c.logger, ErrInvalidOrganizationData)
		return
	}

	response, err := c.service.Create(request)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	responses.DetailResponse(
		ctx,
		http.StatusCreated,
		responses.DetailResponseType[OrganizationResponse]{
			Item:    response,
			Message: "success",
		},
	)
}

// GetByID handles fetching an organization by ID
func (c *Controller) GetByID(ctx *gin.Context) {
	orgID := ctx.Param("id")

	response, err := c.service.GetByID(orgID)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	responses.DetailResponse(
		ctx,
		http.StatusOK,
		responses.DetailResponseType[OrganizationResponse]{
			Item:    response,
			Message: "success",
		},
	)
}

// Update handles updating an organization
func (c *Controller) Update(ctx *gin.Context) {
	orgID := ctx.Param("id")

	var request UpdateOrganizationRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		responses.HandleError(ctx, c.logger, ErrInvalidOrganizationData)
		return
	}

	response, err := c.service.Update(orgID, request)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	responses.DetailResponse(
		ctx,
		http.StatusOK,
		responses.DetailResponseType[OrganizationResponse]{
			Item:    response,
			Message: "success",
		},
	)
}

// List handles fetching a paginated list of organizations
func (c *Controller) List(ctx *gin.Context) {
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

	organizations, total, err := c.service.List(page, limit)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	items := make([]OrganizationListItem, len(organizations))
	for i, org := range organizations {
		items[i] = OrganizationListItem{
			ID:   org.ID.String(),
			Name: org.Name,
		}
	}

	response := OrganizationListResponse{
		Items: items,
		Pagination: responses.PaginationResponseType{
			Total:   total,
			HasNext: (int64(page*limit) < total),
		},
	}

	responses.ListResponse(
		ctx,
		http.StatusOK,
		response,
	)
}
