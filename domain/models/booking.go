package models

import (
	"clean-architecture/pkg/types"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Booking model represents reservations of resources for specific time periods
type Booking struct {
	gorm.Model
	UUID       types.BinaryUUID `json:"uuid" gorm:"index;notnull;unique"`
	ResourceID types.BinaryUUID `json:"resource_id" gorm:"index;not null"`
	UserID     types.BinaryUUID `json:"user_id" gorm:"index;not null"`
	StartTime  time.Time        `json:"start_time" gorm:"not null;index"`
	EndTime    time.Time        `json:"end_time" gorm:"not null;index"`
	Status     string           `json:"status" gorm:"size:50;default:'pending'"`
	Notes      string           `json:"notes" gorm:"type:text"`
	Reference  string           `json:"reference" gorm:"size:100"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (b *Booking) BeforeCreate(tx *gorm.DB) error {
	if b.UUID.String() == (types.BinaryUUID{}).String() {
		id, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		b.UUID = types.BinaryUUID(id)
	}
	return nil
}
