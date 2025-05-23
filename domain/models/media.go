package models

import (
	"gorm.io/gorm"
)

// MediaType represents the type of media file
type MediaType string

const (
	MediaTypeImage    MediaType = "image"
	MediaTypeDocument MediaType = "document"
	MediaTypeVideo    MediaType = "video"
	MediaTypeAudio    MediaType = "audio"
	MediaTypeOther    MediaType = "other"
)

// Media represents a media file in the CMS
type Media struct {
	gorm.Model
	FileName       string    `json:"file_name" gorm:"size:255;not null"`
	FilePath       string    `json:"file_path" gorm:"size:500;not null;uniqueIndex"`
	FileType       string    `json:"file_type" gorm:"size:100;not null"` // MIME type
	FileSize       int64     `json:"file_size" gorm:"not null"`          // Size in bytes
	MediaType      MediaType `json:"media_type" gorm:"size:50;not null"`
	AltText        string    `json:"alt_text" gorm:"size:255"`
	Description    string    `json:"description" gorm:"size:500"`
	Width          int       `json:"width"`    // For images/video
	Height         int       `json:"height"`   // For images/video
	Duration       int       `json:"duration"` // For video/audio in seconds
	UploadedByID   uint      `json:"uploaded_by_id" gorm:"not null"`
	StorageBackend string    `json:"storage_backend" gorm:"size:50;not null;default:'local'"` // local, s3, etc
}
