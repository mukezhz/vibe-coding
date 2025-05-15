package todo

import (
	"clean-architecture/domain/models"
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/types"
)

// Service service layer
type Service struct {
	logger     framework.Logger
	repository Repository
}

// NewService creates a new todo service
func NewService(
	logger framework.Logger,
	repository Repository,
) *Service {
	return &Service{
		logger:     logger,
		repository: repository,
	}
}

// Create creates a new todo
func (s Service) Create(todo *models.Todo) error {
	return s.repository.Create(todo)
}

// GetByID gets a todo by ID
func (s Service) GetByID(todoID types.BinaryUUID) (models.Todo, error) {
	todo, err := s.repository.GetByID(todoID)
	if err != nil {
		// Check if it's a "not found" error
		if err.Error() == "record not found" {
			return todo, NewTodoNotFoundError()
		}
		return todo, err
	}
	return todo, nil
}

// Update updates a todo
func (s Service) Update(todo *models.Todo) error {
	return s.repository.Update(todo)
}

// List returns todos with pagination
func (s Service) List(page, limit int) ([]models.Todo, int64, error) {
	return s.repository.List(page, limit)
}
