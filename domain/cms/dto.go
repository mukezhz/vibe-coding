package cms

import (
	"clean-architecture/domain/models"
	"time"
)

// Content DTOs
type CreateContentRequest struct {
	Title       string             `json:"title" binding:"required"`
	Slug        string             `json:"slug" binding:"required"`
	Body        string             `json:"body"`
	Excerpt     string             `json:"excerpt"`
	FeaturedImg string             `json:"featured_img"`
	Type        models.ContentType `json:"type" binding:"required"`
	CategoryIDs []uint             `json:"category_ids"`
	TagIDs      []uint             `json:"tag_ids"`
}

type UpdateContentRequest struct {
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Body        string `json:"body"`
	Excerpt     string `json:"excerpt"`
	FeaturedImg string `json:"featured_img"`
	CategoryIDs []uint `json:"category_ids"`
	TagIDs      []uint `json:"tag_ids"`
}

type ContentResponse struct {
	ID          uint                 `json:"id"`
	Title       string               `json:"title"`
	Slug        string               `json:"slug"`
	Body        string               `json:"body"`
	Excerpt     string               `json:"excerpt"`
	FeaturedImg string               `json:"featured_img"`
	Type        models.ContentType   `json:"type"`
	Status      models.ContentStatus `json:"status"`
	PublishedAt *time.Time           `json:"published_at"`
	AuthorID    uint                 `json:"author_id"`
	Categories  []CategoryResponse   `json:"categories"`
	Tags        []TagResponse        `json:"tags"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

type ContentListResponse struct {
	ID          uint                 `json:"id"`
	Title       string               `json:"title"`
	Slug        string               `json:"slug"`
	Excerpt     string               `json:"excerpt"`
	Type        models.ContentType   `json:"type"`
	Status      models.ContentStatus `json:"status"`
	PublishedAt *time.Time           `json:"published_at"`
	AuthorID    uint                 `json:"author_id"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

type ContentListQuery struct {
	Page     int                  `form:"page" binding:"min=1"`
	PageSize int                  `form:"page_size" binding:"min=1,max=100"`
	Status   models.ContentStatus `form:"status"`
	Type     models.ContentType   `form:"type"`
	Search   string               `form:"search"`
	SortBy   string               `form:"sort_by"`
	SortDir  string               `form:"sort_dir"`
}

// Category DTOs
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type CategoryResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Tag DTOs
type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

type UpdateTagRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type TagResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Revision DTOs
type RevisionResponse struct {
	ID          uint                 `json:"id"`
	ContentID   uint                 `json:"content_id"`
	Title       string               `json:"title"`
	Body        string               `json:"body"`
	Excerpt     string               `json:"excerpt"`
	Status      models.ContentStatus `json:"status"`
	VersionNum  int                  `json:"version_num"`
	ChangedByID uint                 `json:"changed_by_id"`
	CreatedAt   time.Time            `json:"created_at"`
}

// Media DTOs
type UploadMediaRequest struct {
	AltText     string           `form:"alt_text"`
	Description string           `form:"description"`
	MediaType   models.MediaType `form:"media_type" binding:"required"`
}

type MediaResponse struct {
	ID           uint             `json:"id"`
	FileName     string           `json:"file_name"`
	FileType     string           `json:"file_type"`
	FileSize     int64            `json:"file_size"`
	MediaType    models.MediaType `json:"media_type"`
	AltText      string           `json:"alt_text"`
	Description  string           `json:"description"`
	Width        int              `json:"width,omitempty"`
	Height       int              `json:"height,omitempty"`
	Duration     int              `json:"duration,omitempty"`
	URL          string           `json:"url"`
	ThumbnailURL string           `json:"thumbnail_url,omitempty"`
	UploadedByID uint             `json:"uploaded_by_id"`
	CreatedAt    time.Time        `json:"created_at"`
}

type MediaListQuery struct {
	Page      int              `form:"page" binding:"min=1"`
	PageSize  int              `form:"page_size" binding:"min=1,max=100"`
	MediaType models.MediaType `form:"media_type"`
	Search    string           `form:"search"`
	SortBy    string           `form:"sort_by"`
	SortDir   string           `form:"sort_dir"`
}

// Role DTOs
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RoleResponse struct {
	ID          uint                 `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Permissions []PermissionResponse `json:"permissions,omitempty"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

// Permission DTOs
type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Resource    string `json:"resource" binding:"required"`
	Action      string `json:"action" binding:"required"`
}

type UpdatePermissionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
}

type PermissionResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// User Role DTOs
type AssignRoleRequest struct {
	RoleID uint `json:"role_id" binding:"required"`
}

type AssignPermissionRequest struct {
	PermissionID uint `json:"permission_id" binding:"required"`
}

// Pagination response
type PaginationResponse struct {
	Total       int64 `json:"total"`
	PerPage     int   `json:"per_page"`
	CurrentPage int   `json:"current_page"`
	LastPage    int   `json:"last_page"`
}
