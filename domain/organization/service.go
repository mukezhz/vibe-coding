package organization

import (
	"clean-architecture/domain/models"
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/types"
	"time"
)

// Service service layer
type Service struct {
	repo   *Repository
	logger framework.Logger
}

// NewService creates a new organization service
func NewService(repo *Repository, logger framework.Logger) *Service {
	return &Service{repo, logger}
}

// Create creates a new organization
func (s *Service) Create(request CreateOrganizationRequest) (OrganizationResponse, error) {
	s.logger.Info("[OrganizationService...Create]")

	establishedAt, err := time.Parse("2006-01-02", request.EstablishedAt)
	if err != nil {
		establishedAt = time.Now()
	}

	// Create organization model
	org := models.Organization{
		Name:          request.Name,
		Location:      request.Location,
		EstablishedAt: establishedAt,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Save to database
	if err := s.repo.Create(&org); err != nil {
		return OrganizationResponse{}, err
	}

	// Map to response
	return OrganizationResponse{
		ID:            org.ID.String(),
		Name:          org.Name,
		Location:      org.Location,
		EstablishedAt: org.EstablishedAt,
		CreatedAt:     org.CreatedAt,
		UpdatedAt:     org.UpdatedAt,
	}, nil
}

// GetByID fetches an organization by ID
func (s *Service) GetByID(orgID string) (OrganizationResponse, error) {
	s.logger.Info("[OrganizationService...GetByID]")

	// Convert string ID to BinaryUUID
	id, err := types.ShouldParseUUID(orgID)
	if err != nil {
		return OrganizationResponse{}, ErrInvalidOrganizationData
	}

	// Get from database
	org, err := s.repo.GetByID(id)
	if err != nil {
		return OrganizationResponse{}, ErrOrganizationNotFound
	}

	// Map to response
	return OrganizationResponse{
		ID:            org.ID.String(),
		Name:          org.Name,
		Location:      org.Location,
		EstablishedAt: org.EstablishedAt,
		CreatedAt:     org.CreatedAt,
		UpdatedAt:     org.UpdatedAt,
	}, nil
}

// Update updates an organization
func (s *Service) Update(orgID string, request UpdateOrganizationRequest) (OrganizationResponse, error) {
	s.logger.Info("[OrganizationService...Update]")

	// Convert string ID to BinaryUUID
	id, err := types.ShouldParseUUID(orgID)
	if err != nil {
		return OrganizationResponse{}, ErrInvalidOrganizationData
	}

	// Get existing organization
	org, err := s.repo.GetByID(id)
	if err != nil {
		return OrganizationResponse{}, ErrOrganizationNotFound
	}

	// Update fields if provided
	if request.Name != nil {
		org.Name = *request.Name
	}
	if request.Location != nil {
		org.Location = *request.Location
	}
	if request.EstablishedAt != nil {
		establishedAt, err := time.Parse("2006-01-02", *request.EstablishedAt)
		if err == nil {
			org.EstablishedAt = establishedAt
		}
	}

	org.UpdatedAt = time.Now()

	// Save to database
	if err := s.repo.Update(&org); err != nil {
		return OrganizationResponse{}, err
	}

	// Map to response
	return OrganizationResponse{
		ID:            org.ID.String(),
		Name:          org.Name,
		Location:      org.Location,
		EstablishedAt: org.EstablishedAt,
		CreatedAt:     org.CreatedAt,
		UpdatedAt:     org.UpdatedAt,
	}, nil
}

// List returns a paginated list of organizations
func (s *Service) List(page, limit int) ([]models.Organization, int64, error) {
	s.logger.Info("[OrganizationService...List]")

	// Validate page and limit
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10 // Default limit
	}

	// Get from database with pagination
	orgs, total, err := s.repo.List(page, limit)
	if err != nil {
		return nil, 0, err
	}

	return orgs, total, nil
}
