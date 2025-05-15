package todo

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

// NewRepository creates a new todo repository
func NewRepository(db infrastructure.Database, logger framework.Logger) Repository {
	return Repository{db, logger}
}

// Create creates a new todo
func (r *Repository) Create(todo *models.Todo) error {
	r.logger.Info("[TodoRepository...Create]")
	return r.DB.Create(todo).Error
}

// GetByID gets a todo by ID
func (r *Repository) GetByID(todoID types.BinaryUUID) (todo models.Todo, err error) {
	r.logger.Info("[TodoRepository...GetByID]")
	return todo, r.DB.Where("id = ?", todoID).First(&todo).Error
}

// Update updates a todo
func (r *Repository) Update(todo *models.Todo) error {
	r.logger.Info("[TodoRepository...Update]")
	return r.DB.Save(todo).Error
}

// List returns todos with pagination
func (r *Repository) List(page, limit int) (todos []models.Todo, total int64, err error) {
	r.logger.Info("[TodoRepository...List]")

	offset := (page - 1) * limit

	// Get total count
	if err = r.DB.Model(&models.Todo{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get todos with pagination
	err = r.DB.Offset(offset).Limit(limit).Find(&todos).Error
	return todos, total, err
}
