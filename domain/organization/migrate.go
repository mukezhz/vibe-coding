package organization

import (
	"clean-architecture/domain/models"
	"clean-architecture/pkg/infrastructure"
)

// Migrate automigrates the organization model
func Migrate(db infrastructure.Database) {
	db.AutoMigrate(&models.Organization{})
}
