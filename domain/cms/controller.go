package cms

import (
	"net/http"
	"strconv"

	"clean-architecture/domain/models"
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/responses"
	"clean-architecture/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Controller data type
type Controller struct {
	service *Service
	logger  framework.Logger
	env     *framework.Env
}

// NewController creates new cms controller
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

// Content Controllers

// CreateContent creates a new content item
func (c *Controller) CreateContent(ctx *gin.Context) {
	c.logger.Info("[CMSController...CreateContent]")

	var req CreateContentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	// Get user ID from context (assuming authentication middleware sets this)
	userID := uint(1) // Placeholder for now, should come from authenticated user

	content, err := c.service.CreateContent(&req, userID)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Create response DTO
	response := ContentResponse{
		ID:          content.ID,
		Title:       content.Title,
		Slug:        content.Slug,
		Body:        content.Body,
		Excerpt:     content.Excerpt,
		Type:        content.Type,
		Status:      content.Status,
		AuthorID:    content.AuthorID,
		PublishedAt: content.PublishedAt,
		CreatedAt:   content.CreatedAt,
		UpdatedAt:   content.UpdatedAt,
	}

	responses.DetailResponse(
		ctx,
		http.StatusCreated,
		responses.DetailResponseType[ContentResponse]{
			Item:    response,
			Message: "Content created successfully",
		},
	)
}

func (c *Controller) GetContent(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	content, err := c.service.GetContent(uint(id))
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	responses.DetailResponse(
		ctx,
		http.StatusOK,
		responses.DetailResponseType[models.Content]{
			Item:    utils.SafeDeref(content),
			Message: "Content retrieved successfully",
		},
	)
}

func (c *Controller) UpdateContent(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	var req UpdateContentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	userID := uint(1) // Placeholder for now, should come from authenticated user

	content, err := c.service.UpdateContent(uint(id), &req, userID)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	responses.DetailResponse(
		ctx,
		http.StatusOK,
		responses.DetailResponseType[models.Content]{
			Item:    utils.SafeDeref(content),
			Message: "Content updated successfully",
		},
	)
}

func (c *Controller) DeleteContent(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	err = c.service.DeleteContent(uint(id))
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Content deleted successfully",
	})
}

func (c *Controller) ListContent(ctx *gin.Context) {
	var query ContentListQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	// Set default values if not provided
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 100 {
		query.PageSize = 10
	}

	contents, pagination, err := c.service.ListContent(&query)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Convert to response DTOs
	var contentResponses []ContentListResponse
	for _, content := range contents {
		contentResponses = append(contentResponses, ContentListResponse{
			ID:          content.ID,
			Title:       content.Title,
			Slug:        content.Slug,
			Excerpt:     content.Excerpt,
			Type:        content.Type,
			Status:      content.Status,
			PublishedAt: content.PublishedAt,
			AuthorID:    content.AuthorID,
			CreatedAt:   content.CreatedAt,
			UpdatedAt:   content.UpdatedAt,
		})
	}

	response := map[string]interface{}{
		"data":       contentResponses,
		"pagination": pagination,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Content list retrieved successfully",
		"data":    response,
	})
}

func (c *Controller) PublishContent(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	// Get user ID from context (assuming authentication middleware sets this)
	userID := uint(1) // Placeholder for now, should come from authenticated user

	content, err := c.service.PublishContent(uint(id), userID)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Content published successfully",
		"data":    content,
	})
}

func (c *Controller) UnpublishContent(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	// Get user ID from context (assuming authentication middleware sets this)
	userID := uint(1) // Placeholder for now, should come from authenticated user

	content, err := c.service.UnpublishContent(uint(id), userID)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Content unpublished successfully",
		"data":    content,
	})
}

func (c *Controller) GetContentRevisions(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	revisions, err := c.service.GetContentRevisions(uint(id))
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Convert to response DTOs
	var revisionResponses []RevisionResponse
	for _, revision := range revisions {
		revisionResponses = append(revisionResponses, RevisionResponse{
			ID:          revision.ID,
			ContentID:   revision.ContentID,
			Title:       revision.Title,
			Body:        revision.Body,
			Excerpt:     revision.Excerpt,
			Status:      revision.Status,
			VersionNum:  revision.VersionNum,
			ChangedByID: revision.ChangedByID,
			CreatedAt:   revision.CreatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Content revisions retrieved successfully",
		"data":    revisionResponses,
	})
}

func (c *Controller) GetContentRevision(ctx *gin.Context) {
	contentID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	revisionID, err := strconv.ParseUint(ctx.Param("revisionId"), 10, 32)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	revision, err := c.service.GetContentRevision(uint(contentID), uint(revisionID))
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	revisionResponse := RevisionResponse{
		ID:          revision.ID,
		ContentID:   revision.ContentID,
		Title:       revision.Title,
		Body:        revision.Body,
		Excerpt:     revision.Excerpt,
		Status:      revision.Status,
		VersionNum:  revision.VersionNum,
		ChangedByID: revision.ChangedByID,
		CreatedAt:   revision.CreatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Content revision retrieved successfully",
		"data":    revisionResponse,
	})
}

// Category Controllers
func (c *Controller) CreateCategory(ctx *gin.Context) {
	var req CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	category, err := c.service.CreateCategory(&req)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Category created successfully",
		"data":    category,
	})
}

func (c *Controller) GetCategory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	category, err := c.service.GetCategory(uint(id))
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Category retrieved successfully",
		"data":    category,
	})
}

func (c *Controller) UpdateCategory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	var req UpdateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	category, err := c.service.UpdateCategory(uint(id), &req)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Category updated successfully",
		"data":    category,
	})
}

func (c *Controller) DeleteCategory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	err = c.service.DeleteCategory(uint(id))
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Category deleted successfully",
	})
}

func (c *Controller) ListCategories(ctx *gin.Context) {
	categories, err := c.service.ListCategories()
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Convert to response DTOs
	var categoryResponses []CategoryResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Slug:        category.Slug,
			Description: category.Description,
			CreatedAt:   category.CreatedAt,
			UpdatedAt:   category.UpdatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Categories retrieved successfully",
		"data":    categoryResponses,
	})
}

// Tag Controllers
func (c *Controller) CreateTag(ctx *gin.Context) {
	var req CreateTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	tag, err := c.service.CreateTag(&req)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Tag created successfully",
		"data":    tag,
	})
}

func (c *Controller) GetTag(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	tag, err := c.service.GetTag(uint(id))
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Tag retrieved successfully",
		"data":    tag,
	})
}

func (c *Controller) UpdateTag(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	var req UpdateTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	tag, err := c.service.UpdateTag(uint(id), &req)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Tag updated successfully",
		"data":    tag,
	})
}

func (c *Controller) DeleteTag(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	err = c.service.DeleteTag(uint(id))
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Tag deleted successfully",
	})
}

func (c *Controller) ListTags(ctx *gin.Context) {
	tags, err := c.service.ListTags()
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Convert to response DTOs
	var tagResponses []TagResponse
	for _, tag := range tags {
		tagResponses = append(tagResponses, TagResponse{
			ID:        tag.ID,
			Name:      tag.Name,
			Slug:      tag.Slug,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Tags retrieved successfully",
		"data":    tagResponses,
	})
}

// Media Controllers
func (c *Controller) UploadMedia(ctx *gin.Context) {
	// Get file from request
	file, err := ctx.FormFile("file")
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	var req UploadMediaRequest
	if err := ctx.ShouldBind(&req); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	// Get user ID from context (assuming authentication middleware sets this)
	userID := uint(1) // Placeholder for now, should come from authenticated user

	media, err := c.service.UploadMedia(file, &req, userID)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Construct response with URLs
	mediaResponse := MediaResponse{
		ID:           media.ID,
		FileName:     media.FileName,
		FileType:     media.FileType,
		FileSize:     media.FileSize,
		MediaType:    media.MediaType,
		AltText:      media.AltText,
		Description:  media.Description,
		Width:        media.Width,
		Height:       media.Height,
		Duration:     media.Duration,
		URL:          media.FilePath, // In production, this should be a full URL
		UploadedByID: media.UploadedByID,
		CreatedAt:    media.CreatedAt,
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Media uploaded successfully",
		"data":    mediaResponse,
	})
}

func (c *Controller) GetMedia(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	media, err := c.service.GetMedia(uint(id))
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Construct response with URLs
	mediaResponse := MediaResponse{
		ID:           media.ID,
		FileName:     media.FileName,
		FileType:     media.FileType,
		FileSize:     media.FileSize,
		MediaType:    media.MediaType,
		AltText:      media.AltText,
		Description:  media.Description,
		Width:        media.Width,
		Height:       media.Height,
		Duration:     media.Duration,
		URL:          media.FilePath, // In production, this should be a full URL
		UploadedByID: media.UploadedByID,
		CreatedAt:    media.CreatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Media retrieved successfully",
		"data":    mediaResponse,
	})
}

func (c *Controller) DeleteMedia(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	err = c.service.DeleteMedia(uint(id))
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Media deleted successfully",
	})
}

func (c *Controller) ListMedia(ctx *gin.Context) {
	var query MediaListQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	// Set default values if not provided
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 100 {
		query.PageSize = 10
	}

	media, pagination, err := c.service.ListMedia(&query)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Convert to response DTOs
	var mediaResponses []MediaResponse
	for _, m := range media {
		mediaResponses = append(mediaResponses, MediaResponse{
			ID:           m.ID,
			FileName:     m.FileName,
			FileType:     m.FileType,
			FileSize:     m.FileSize,
			MediaType:    m.MediaType,
			AltText:      m.AltText,
			Description:  m.Description,
			Width:        m.Width,
			Height:       m.Height,
			Duration:     m.Duration,
			URL:          m.FilePath, // In production, this should be a full URL
			UploadedByID: m.UploadedByID,
			CreatedAt:    m.CreatedAt,
		})
	}

	response := map[string]interface{}{
		"data":       mediaResponses,
		"pagination": pagination,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Media list retrieved successfully",
		"data":    response,
	})
}

// Role Controllers
func (c *Controller) CreateRole(ctx *gin.Context) {
	var req CreateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	role, err := c.service.CreateRole(&req)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Role created successfully",
		"data":    role,
	})
}

func (c *Controller) GetRole(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	role, err := c.service.GetRole(uint(id))
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Convert to response DTO
	var permissionResponses []PermissionResponse
	for _, permission := range role.Permissions {
		permissionResponses = append(permissionResponses, PermissionResponse{
			ID:          permission.ID,
			Name:        permission.Name,
			Description: permission.Description,
			Resource:    permission.Resource,
			Action:      permission.Action,
			CreatedAt:   permission.CreatedAt,
			UpdatedAt:   permission.UpdatedAt,
		})
	}

	roleResponse := RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Permissions: permissionResponses,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Role retrieved successfully",
		"data":    roleResponse,
	})
}

func (c *Controller) UpdateRole(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	var req UpdateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	role, err := c.service.UpdateRole(uint(id), &req)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Role updated successfully",
		"data":    role,
	})
}

func (c *Controller) DeleteRole(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	err = c.service.DeleteRole(uint(id))
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Role deleted successfully",
	})
}

func (c *Controller) ListRoles(ctx *gin.Context) {
	roles, err := c.service.ListRoles()
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Convert to response DTOs
	var roleResponses []RoleResponse
	for _, role := range roles {
		var permissionResponses []PermissionResponse
		for _, permission := range role.Permissions {
			permissionResponses = append(permissionResponses, PermissionResponse{
				ID:          permission.ID,
				Name:        permission.Name,
				Description: permission.Description,
				Resource:    permission.Resource,
				Action:      permission.Action,
				CreatedAt:   permission.CreatedAt,
				UpdatedAt:   permission.UpdatedAt,
			})
		}

		roleResponses = append(roleResponses, RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Permissions: permissionResponses,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Roles retrieved successfully",
		"data":    roleResponses,
	})
}

func (c *Controller) AssignPermissionToRole(ctx *gin.Context) {
	roleID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	var req AssignPermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	err = c.service.AssignPermissionToRole(uint(roleID), req.PermissionID)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Permission assigned to role successfully",
	})
}

func (c *Controller) RemovePermissionFromRole(ctx *gin.Context) {
	roleID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	permissionID, err := strconv.ParseUint(ctx.Param("permissionId"), 10, 32)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	err = c.service.RemovePermissionFromRole(uint(roleID), uint(permissionID))
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Permission removed from role successfully",
	})
}

// Permission Controllers
func (c *Controller) CreatePermission(ctx *gin.Context) {
	var req CreatePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	permission, err := c.service.CreatePermission(&req)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Permission created successfully",
		"data":    permission,
	})
}

func (c *Controller) GetPermission(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	permission, err := c.service.GetPermission(uint(id))
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Permission retrieved successfully",
		"data":    permission,
	})
}

func (c *Controller) UpdatePermission(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	var req UpdatePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	permission, err := c.service.UpdatePermission(uint(id), &req)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Permission updated successfully",
		"data":    permission,
	})
}

func (c *Controller) DeletePermission(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	err = c.service.DeletePermission(uint(id))
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Permission deleted successfully",
	})
}

func (c *Controller) ListPermissions(ctx *gin.Context) {
	permissions, err := c.service.ListPermissions()
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Convert to response DTOs
	var permissionResponses []PermissionResponse
	for _, permission := range permissions {
		permissionResponses = append(permissionResponses, PermissionResponse{
			ID:          permission.ID,
			Name:        permission.Name,
			Description: permission.Description,
			Resource:    permission.Resource,
			Action:      permission.Action,
			CreatedAt:   permission.CreatedAt,
			UpdatedAt:   permission.UpdatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Permissions retrieved successfully",
		"data":    permissionResponses,
	})
}

// User Role Controllers
func (c *Controller) AssignRoleToUser(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	var req AssignRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	err = c.service.AssignRoleToUser(uint(userID), req.RoleID)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	responses.DetailResponse(ctx, http.StatusOK, responses.DetailResponseType[any]{
		Message: "Role assigned to user successfully",
	})
}

func (c *Controller) RemoveRoleFromUser(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	roleID, err := strconv.ParseUint(ctx.Param("roleId"), 10, 32)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	err = c.service.RemoveRoleFromUser(uint(userID), uint(roleID))
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	responses.DetailResponse(ctx, http.StatusOK, responses.DetailResponseType[any]{
		Message: "Role removed from user successfully",
	})
}

func (c *Controller) GetUserRoles(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	roles, err := c.service.GetUserRoles(uint(userID))
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Convert to response DTOs
	var roleResponses []RoleResponse
	for _, role := range roles {
		var permissionResponses []PermissionResponse
		for _, permission := range role.Permissions {
			permissionResponses = append(permissionResponses, PermissionResponse{
				ID:          permission.ID,
				Name:        permission.Name,
				Description: permission.Description,
				Resource:    permission.Resource,
				Action:      permission.Action,
				CreatedAt:   permission.CreatedAt,
				UpdatedAt:   permission.UpdatedAt,
			})
		}

		roleResponses = append(roleResponses, RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Permissions: permissionResponses,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		})
	}

	responses.DetailResponse(ctx, http.StatusOK, responses.DetailResponseType[any]{
		Message: "User roles retrieved successfully",
	})
}
