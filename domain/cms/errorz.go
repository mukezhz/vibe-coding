package cms

import (
	"clean-architecture/pkg/errorz"
)

var (
	// Content errors
	ErrContentNotFound     = errorz.ErrNotFound.JoinError("Content not found")
	ErrContentCreateFailed = errorz.ErrInternal.JoinError("Failed to create content")
	ErrContentUpdateFailed = errorz.ErrInternal.JoinError("Failed to update content")
	ErrContentDeleteFailed = errorz.ErrInternal.JoinError("Failed to delete content")
	ErrSlugAlreadyExists   = errorz.ErrBadRequest.JoinError("Content with this slug already exists")

	// Category errors
	ErrCategoryNotFound     = errorz.ErrNotFound.JoinError("Category not found")
	ErrCategoryCreateFailed = errorz.ErrInternal.JoinError("Failed to create category")
	ErrCategoryUpdateFailed = errorz.ErrInternal.JoinError("Failed to update category")
	ErrCategoryDeleteFailed = errorz.ErrInternal.JoinError("Failed to delete category")
	ErrCategorySlugExists   = errorz.ErrBadRequest.JoinError("Category with this slug already exists")

	// Tag errors
	ErrTagNotFound     = errorz.ErrNotFound.JoinError("Tag not found")
	ErrTagCreateFailed = errorz.ErrInternal.JoinError("Failed to create tag")
	ErrTagUpdateFailed = errorz.ErrInternal.JoinError("Failed to update tag")
	ErrTagDeleteFailed = errorz.ErrInternal.JoinError("Failed to delete tag")
	ErrTagSlugExists   = errorz.ErrBadRequest.JoinError("Tag with this slug already exists")

	// Revision errors
	ErrRevisionNotFound = errorz.ErrNotFound.JoinError("Revision not found")

	// Media errors
	ErrMediaNotFound     = errorz.ErrNotFound.JoinError("Media not found")
	ErrMediaCreateFailed = errorz.ErrInternal.JoinError("Failed to upload media")
	ErrMediaDeleteFailed = errorz.ErrInternal.JoinError("Failed to delete media")
	ErrInvalidMediaType  = errorz.ErrBadRequest.JoinError("Invalid media type")
	ErrFileTooLarge      = errorz.ErrBadRequest.JoinError("File size exceeds the maximum limit")

	// Role errors
	ErrRoleNotFound     = errorz.ErrNotFound.JoinError("Role not found")
	ErrRoleCreateFailed = errorz.ErrInternal.JoinError("Failed to create role")
	ErrRoleUpdateFailed = errorz.ErrInternal.JoinError("Failed to update role")
	ErrRoleDeleteFailed = errorz.ErrInternal.JoinError("Failed to delete role")
	ErrRoleNameExists   = errorz.ErrBadRequest.JoinError("Role with this name already exists")

	// Permission errors
	ErrPermissionNotFound     = errorz.ErrNotFound.JoinError("Permission not found")
	ErrPermissionCreateFailed = errorz.ErrInternal.JoinError("Failed to create permission")
	ErrPermissionUpdateFailed = errorz.ErrInternal.JoinError("Failed to update permission")
	ErrPermissionDeleteFailed = errorz.ErrInternal.JoinError("Failed to delete permission")
	ErrPermissionNameExists   = errorz.ErrBadRequest.JoinError("Permission with this name already exists")

	// User Role errors
	ErrUserRoleCreateFailed = errorz.ErrInternal.JoinError("Failed to assign role to user")
	ErrUserRoleDeleteFailed = errorz.ErrInternal.JoinError("Failed to remove role from user")
	ErrRolePermissionFailed = errorz.ErrInternal.JoinError("Failed to manage permission for role")

	// Authorization errors
	ErrUnauthorized = errorz.ErrUnauthorized.JoinError("You are not authorized to perform this action")
	ErrForbidden    = errorz.ErrForbidden.JoinError("You do not have permission to perform this action")
)
