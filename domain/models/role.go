package models

import (
	"gorm.io/gorm"
)

// Role represents a user role in the system
type Role struct {
	gorm.Model
	Name        string       `json:"name" gorm:"size:50;not null;uniqueIndex"`
	Description string       `json:"description" gorm:"size:255"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
}

// Permission represents a specific permission in the system
type Permission struct {
	gorm.Model
	Name        string `json:"name" gorm:"size:100;not null;uniqueIndex"`
	Description string `json:"description" gorm:"size:255"`
	Resource    string `json:"resource" gorm:"size:50;not null"` // e.g., "content", "user", "media"
	Action      string `json:"action" gorm:"size:50;not null"`   // e.g., "create", "read", "update", "delete"
	Roles       []Role `json:"roles" gorm:"many2many:role_permissions;"`
}

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	ID     uint `gorm:"primarykey"`
	UserID uint `json:"user_id" gorm:"not null"`
	RoleID uint `json:"role_id" gorm:"not null"`
}
