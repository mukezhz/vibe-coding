package cms

import (
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/infrastructure"
)

// Route struct for CMS
type Route struct {
	logger     framework.Logger
	handler    infrastructure.Router
	controller *Controller
}

// NewRoute creates a new route
func NewRoute(
	logger framework.Logger,
	handler infrastructure.Router,
	controller *Controller,
) *Route {
	return &Route{
		handler:    handler,
		logger:     logger,
		controller: controller,
	}
}

// RegisterRoute registers the CMS routes
func RegisterRoute(r *Route) {
	r.logger.Info("Setting up CMS routes")

	// Group all API routes under /api/v1/cms
	api := r.handler.Group("/api/v1/cms")

	// Content routes
	content := api.Group("/content")
	{
		content.POST("", r.controller.CreateContent)
		content.GET("", r.controller.ListContent)
		content.GET("/:id", r.controller.GetContent)
		content.PUT("/:id", r.controller.UpdateContent)
		content.DELETE("/:id", r.controller.DeleteContent)
		content.PUT("/:id/publish", r.controller.PublishContent)
		content.PUT("/:id/unpublish", r.controller.UnpublishContent)
		content.GET("/:id/revisions", r.controller.GetContentRevisions)
		content.GET("/:id/revisions/:revisionId", r.controller.GetContentRevision)
	}

	// Category routes
	category := api.Group("/categories")
	{
		category.POST("", r.controller.CreateCategory)
		category.GET("", r.controller.ListCategories)
		category.GET("/:id", r.controller.GetCategory)
		category.PUT("/:id", r.controller.UpdateCategory)
		category.DELETE("/:id", r.controller.DeleteCategory)
	}

	// Tag routes
	tag := api.Group("/tags")
	{
		tag.POST("", r.controller.CreateTag)
		tag.GET("", r.controller.ListTags)
		tag.GET("/:id", r.controller.GetTag)
		tag.PUT("/:id", r.controller.UpdateTag)
		tag.DELETE("/:id", r.controller.DeleteTag)
	}

	// Media routes
	media := api.Group("/media")
	{
		media.POST("", r.controller.UploadMedia)
		media.GET("", r.controller.ListMedia)
		media.GET("/:id", r.controller.GetMedia)
		media.DELETE("/:id", r.controller.DeleteMedia)
	}

	// Role routes
	role := api.Group("/roles")
	{
		role.POST("", r.controller.CreateRole)
		role.GET("", r.controller.ListRoles)
		role.GET("/:id", r.controller.GetRole)
		role.PUT("/:id", r.controller.UpdateRole)
		role.DELETE("/:id", r.controller.DeleteRole)
		role.POST("/:id/permissions", r.controller.AssignPermissionToRole)
		role.DELETE("/:id/permissions/:permissionId", r.controller.RemovePermissionFromRole)
	}

	// Permission routes
	permission := api.Group("/permissions")
	{
		permission.POST("", r.controller.CreatePermission)
		permission.GET("", r.controller.ListPermissions)
		permission.GET("/:id", r.controller.GetPermission)
		permission.PUT("/:id", r.controller.UpdatePermission)
		permission.DELETE("/:id", r.controller.DeletePermission)
	}

	// User role management
	userRole := api.Group("/users")
	{
		userRole.POST("/:id/roles", r.controller.AssignRoleToUser)
		userRole.DELETE("/:id/roles/:roleId", r.controller.RemoveRoleFromUser)
		userRole.GET("/:id/roles", r.controller.GetUserRoles)
	}
}
