package models

import (
	"clean-architecture/pkg/types"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Availability model represents time slots when resources are available for booking
type Availability struct {
	gorm.Model
	UUID        types.BinaryUUID `json:"uuid" gorm:"index;notnull;unique"`
	ResourceID  types.BinaryUUID `json:"resource_id" gorm:"index;not null"`
	StartTime   time.Time        `json:"start_time" gorm:"not null;index"`
	EndTime     time.Time        `json:"end_time" gorm:"not null;index"`
	IsRecurring bool             `json:"is_recurring" gorm:"default:false"`
	RecurRule   string           `json:"recur_rule" gorm:"size:255"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (a *Availability) BeforeCreate(tx *gorm.DB) error {
	if a.UUID.String() == (types.BinaryUUID{}).String() {
		id, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		a.UUID = types.BinaryUUID(id)
	}
	return nil
}
