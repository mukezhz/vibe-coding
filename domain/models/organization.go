package models

import (
	"time"

	"clean-architecture/pkg/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Organization represents the organization model in the database
type Organization struct {
	ID            types.BinaryUUID `json:"id" gorm:"type:binary(16);primary_key"`
	Name          string           `json:"name" gorm:"not null"`
	Location      string           `json:"location"`
	EstablishedAt time.Time        `json:"established_at"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

func (Organization) TableName() string {
	return "organizations"
}

func (u *Organization) BeforeCreate(tx *gorm.DB) error {
	if u.ID.String() == (types.BinaryUUID{}).String() {
		id, err := uuid.NewRandom()
		u.ID = types.BinaryUUID(id)
		return err
	}
	return nil
}
