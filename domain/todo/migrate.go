package todo

import (
	"clean-architecture/domain/models"
	"clean-architecture/pkg/infrastructure"
)

// Migrate automigrates the todo model
func Migrate(db infrastructure.Database) {
	db.AutoMigrate(&models.Todo{})
}
