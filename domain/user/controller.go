package user

import (
	"clean-architecture/domain/models"
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/responses"
	"clean-architecture/pkg/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserController data type
type Controller struct {
	service *Service
	logger  framework.Logger
	env     *framework.Env
}

type URLObject struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// NewUserController creates new user controller
func NewController(
	userService *Service,
	logger framework.Logger,
	env *framework.Env,
) *Controller {
	return &Controller{
		service: userService,
		logger:  logger,
		env:     env,
	}
}

// CreateUser creates the new user
func (c *Controller) CreateUser(ctx *gin.Context) {
	var user models.User

	if err := ctx.Bind(&user); err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// check if the user already exists
	if err := c.service.Create(&user); err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	responses.MessageOnlyResponse(ctx, http.StatusCreated, "User created successfully")
}

// GetOneUser gets one user
func (c *Controller) GetUserByID(ctx *gin.Context) {
	paramID := ctx.Param("id")

	userID, err := types.ShouldParseUUID(paramID)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, ErrInvalidUserID)
		return
	}

	user, err := c.service.GetUserByID(userID)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	responses.DetailResponse(
		ctx,
		http.StatusOK,
		responses.DetailResponseType[models.User]{
			Item:    user,
			Message: "success",
		},
	)
}
