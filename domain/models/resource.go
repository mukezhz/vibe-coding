package models

import (
	"clean-architecture/pkg/types"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Resource model represents bookable items like rooms, equipment, services
type Resource struct {
	gorm.Model
	UUID        types.BinaryUUID `json:"uuid" gorm:"index;notnull;unique"`
	Name        string           `json:"name" gorm:"size:255;not null"`
	Description string           `json:"description" gorm:"type:text"`
	Type        string           `json:"type" gorm:"size:50;not null"`
	Capacity    int              `json:"capacity" gorm:"default:1"`
	Location    string           `json:"location" gorm:"size:255"`
	Attributes  datatypes.JSON   `json:"attributes" gorm:"type:json"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (r *Resource) BeforeCreate(tx *gorm.DB) error {
	if r.UUID.String() == (types.BinaryUUID{}).String() {
		id, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		r.UUID = types.BinaryUUID(id)
	}
	return nil
}
