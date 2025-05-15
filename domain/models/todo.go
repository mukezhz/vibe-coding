package models

import (
	"time"

	"clean-architecture/pkg/types"
)

// Todo represents the todo model in the database
type Todo struct {
	ID          types.BinaryUUID `json:"id" gorm:"type:binary(16);primary_key"`
	Title       string           `json:"title" gorm:"not null"`
	Description string           `json:"description"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}
