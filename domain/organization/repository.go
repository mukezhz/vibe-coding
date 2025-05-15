package organization

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

// NewRepository creates a new organization repository
func NewRepository(db infrastructure.Database, logger framework.Logger) *Repository {
	return &Repository{db, logger}
}

// Create creates a new organization
func (r *Repository) Create(org *models.Organization) error {
	r.logger.Info("[OrganizationRepository...Create]")
	return r.DB.Create(org).Error
}

// GetByID gets an organization by ID
func (r *Repository) GetByID(orgID types.BinaryUUID) (org models.Organization, err error) {
	r.logger.Info("[OrganizationRepository...GetByID]")
	return org, r.DB.Where("id = ?", orgID).First(&org).Error
}

// Update updates an organization
func (r *Repository) Update(org *models.Organization) error {
	r.logger.Info("[OrganizationRepository...Update]")
	return r.DB.Save(org).Error
}

// List returns organizations with pagination
func (r *Repository) List(page, limit int) (orgs []models.Organization, total int64, err error) {
	r.logger.Info("[OrganizationRepository...List]")

	offset := (page - 1) * limit

	// Get total count
	if err = r.DB.Model(&models.Organization{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get organizations with pagination
	err = r.DB.Offset(offset).Limit(limit).Find(&orgs).Error
	return orgs, total, err
}
