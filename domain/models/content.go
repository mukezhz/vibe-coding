package models

import (
	"time"

	"gorm.io/gorm"
)

// ContentStatus represents the publishing status of content
type ContentStatus string

const (
	ContentStatusDraft     ContentStatus = "draft"
	ContentStatusPublished ContentStatus = "published"
	ContentStatusArchived  ContentStatus = "archived"
)

// ContentType represents the type of content
type ContentType string

const (
	ContentTypeArticle ContentType = "article"
	ContentTypePage    ContentType = "page"
)

// Content represents an article, page or other content in the CMS
type Content struct {
	gorm.Model
	Title       string        `json:"title" gorm:"size:255;not null"`
	Slug        string        `json:"slug" gorm:"size:255;not null;uniqueIndex"`
	Body        string        `json:"body" gorm:"type:text"`
	Excerpt     string        `json:"excerpt" gorm:"size:500"`
	FeaturedImg string        `json:"featured_img" gorm:"size:255"`
	Type        ContentType   `json:"type" gorm:"size:50;not null;default:'article'"`
	Status      ContentStatus `json:"status" gorm:"size:50;not null;default:'draft'"`
	PublishedAt *time.Time    `json:"published_at"`
	AuthorID    uint          `json:"author_id" gorm:"not null"`
	Categories  []Category    `json:"categories" gorm:"many2many:content_categories;"`
	Tags        []Tag         `json:"tags" gorm:"many2many:content_tags;"`
	Revisions   []Revision    `json:"revisions" gorm:"foreignKey:ContentID"`
}

// Category for categorizing content
type Category struct {
	gorm.Model
	Name        string    `json:"name" gorm:"size:100;not null;uniqueIndex"`
	Slug        string    `json:"slug" gorm:"size:100;not null;uniqueIndex"`
	Description string    `json:"description" gorm:"size:500"`
	Contents    []Content `json:"contents" gorm:"many2many:content_categories;"`
}

// Tag for tagging content
type Tag struct {
	gorm.Model
	Name     string    `json:"name" gorm:"size:50;not null;uniqueIndex"`
	Slug     string    `json:"slug" gorm:"size:50;not null;uniqueIndex"`
	Contents []Content `json:"contents" gorm:"many2many:content_tags;"`
}

// Revision represents a version of content
type Revision struct {
	gorm.Model
	ContentID   uint          `json:"content_id" gorm:"not null"`
	Title       string        `json:"title" gorm:"size:255;not null"`
	Body        string        `json:"body" gorm:"type:text"`
	Excerpt     string        `json:"excerpt" gorm:"size:500"`
	Status      ContentStatus `json:"status" gorm:"size:50;not null"`
	VersionNum  int           `json:"version_num" gorm:"not null"`
	ChangedByID uint          `json:"changed_by_id" gorm:"not null"`
}
