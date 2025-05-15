package organization

import (
	"clean-architecture/pkg/responses"
	"time"
)

// CreateOrganizationRequest DTO for creating an organization
type CreateOrganizationRequest struct {
	Name          string `json:"name" binding:"required"`
	Location      string `json:"location"`
	EstablishedAt string `json:"established_at"`
}

// OrganizationResponse DTO for organization response
type OrganizationResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Location      string    `json:"location"`
	EstablishedAt time.Time `json:"established_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// OrganizationListItem DTO for items in organization list
type OrganizationListItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// OrganizationListResponse DTO for paginated organization list
type OrganizationListResponse = responses.ListResponseType[OrganizationListItem]

// PageInfo contains pagination information
type PageInfo struct {
	HasNext bool  `json:"has_next"`
	Total   int64 `json:"total"`
}

// UpdateOrganizationRequest DTO for updating an organization
type UpdateOrganizationRequest struct {
	Name          *string `json:"name"`
	Location      *string `json:"location"`
	EstablishedAt *string `json:"established_at"`
}
