package cms

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"clean-architecture/domain/models"
	"clean-architecture/pkg/utils"

	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
	db   *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		repo: NewRepository(db),
		db:   db,
	}
}

// Content Service Methods
func (s *Service) CreateContent(req *CreateContentRequest, userID uint) (*models.Content, error) {
	// Check if slug already exists
	existingContent, err := s.repo.GetContentBySlug(req.Slug)
	if err == nil && existingContent != nil {
		return nil, ErrSlugAlreadyExists
	}

	content := &models.Content{
		Title:       req.Title,
		Slug:        req.Slug,
		Body:        req.Body,
		Excerpt:     req.Excerpt,
		FeaturedImg: req.FeaturedImg,
		Type:        req.Type,
		Status:      models.ContentStatusDraft,
		AuthorID:    userID,
	}

	// Start transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := NewRepository(tx)

		// Create content
		if err := txRepo.CreateContent(content); err != nil {
			return err
		}

		// Add categories if provided
		if len(req.CategoryIDs) > 0 {
			for _, catID := range req.CategoryIDs {
				cat, err := txRepo.GetCategoryByID(catID)
				if err != nil {
					return err
				}
				if err := tx.Model(content).Association("Categories").Append(cat); err != nil {
					return err
				}
			}
		}

		// Add tags if provided
		if len(req.TagIDs) > 0 {
			for _, tagID := range req.TagIDs {
				tag, err := txRepo.GetTagByID(tagID)
				if err != nil {
					return err
				}
				if err := tx.Model(content).Association("Tags").Append(tag); err != nil {
					return err
				}
			}
		}

		// Create initial revision
		revision := &models.Revision{
			ContentID:   content.ID,
			Title:       content.Title,
			Body:        content.Body,
			Excerpt:     content.Excerpt,
			Status:      content.Status,
			VersionNum:  1,
			ChangedByID: userID,
		}
		if err := txRepo.CreateRevision(revision); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, ErrContentCreateFailed
	}

	return content, nil
}

func (s *Service) GetContent(id uint) (*models.Content, error) {
	return s.repo.GetContentByID(id)
}

func (s *Service) UpdateContent(id uint, req *UpdateContentRequest, userID uint) (*models.Content, error) {
	content, err := s.repo.GetContentByID(id)
	if err != nil {
		return nil, err
	}

	// Check if slug is being changed and if it already exists
	if req.Slug != "" && req.Slug != content.Slug {
		existingContent, err := s.repo.GetContentBySlug(req.Slug)
		if err == nil && existingContent != nil && existingContent.ID != id {
			return nil, ErrSlugAlreadyExists
		}
	}

	// Start transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := NewRepository(tx)

		// Update content fields
		if req.Title != "" {
			content.Title = req.Title
		}
		if req.Slug != "" {
			content.Slug = req.Slug
		}
		if req.Body != "" {
			content.Body = req.Body
		}
		if req.Excerpt != "" {
			content.Excerpt = req.Excerpt
		}
		if req.FeaturedImg != "" {
			content.FeaturedImg = req.FeaturedImg
		}

		// Save content
		if err := txRepo.UpdateContent(content); err != nil {
			return err
		}

		// Update categories if provided
		if len(req.CategoryIDs) > 0 {
			// Clear existing categories
			if err := tx.Model(content).Association("Categories").Clear(); err != nil {
				return err
			}

			// Add new categories
			for _, catID := range req.CategoryIDs {
				cat, err := txRepo.GetCategoryByID(catID)
				if err != nil {
					return err
				}
				if err := tx.Model(content).Association("Categories").Append(cat); err != nil {
					return err
				}
			}
		}

		// Update tags if provided
		if len(req.TagIDs) > 0 {
			// Clear existing tags
			if err := tx.Model(content).Association("Tags").Clear(); err != nil {
				return err
			}

			// Add new tags
			for _, tagID := range req.TagIDs {
				tag, err := txRepo.GetTagByID(tagID)
				if err != nil {
					return err
				}
				if err := tx.Model(content).Association("Tags").Append(tag); err != nil {
					return err
				}
			}
		}

		// Create new revision
		latestRevisions, err := txRepo.GetContentRevisions(content.ID)
		if err != nil {
			return err
		}

		var versionNum int
		if len(latestRevisions) > 0 {
			versionNum = latestRevisions[0].VersionNum + 1
		} else {
			versionNum = 1
		}

		revision := &models.Revision{
			ContentID:   content.ID,
			Title:       content.Title,
			Body:        content.Body,
			Excerpt:     content.Excerpt,
			Status:      content.Status,
			VersionNum:  versionNum,
			ChangedByID: userID,
		}
		if err := txRepo.CreateRevision(revision); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, ErrContentUpdateFailed
	}

	// Get updated content
	return s.repo.GetContentByID(id)
}

func (s *Service) DeleteContent(id uint) error {
	_, err := s.repo.GetContentByID(id)
	if err != nil {
		return err
	}

	err = s.repo.DeleteContent(id)
	if err != nil {
		return ErrContentDeleteFailed
	}

	return nil
}

func (s *Service) ListContent(query *ContentListQuery) ([]models.Content, *PaginationResponse, error) {
	// Set default values if not provided
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 100 {
		query.PageSize = 10
	}

	contents, total, err := s.repo.ListContent(query)
	if err != nil {
		return nil, nil, err
	}

	// Calculate pagination info
	totalPages := (int(total) + query.PageSize - 1) / query.PageSize
	pagination := &PaginationResponse{
		Total:       total,
		PerPage:     query.PageSize,
		CurrentPage: query.Page,
		LastPage:    totalPages,
	}

	return contents, pagination, nil
}

func (s *Service) PublishContent(id uint, userID uint) (*models.Content, error) {
	content, err := s.repo.GetContentByID(id)
	if err != nil {
		return nil, err
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := NewRepository(tx)

		// Publish content
		if err := txRepo.PublishContent(content); err != nil {
			return err
		}

		// Create new revision
		latestRevisions, err := txRepo.GetContentRevisions(content.ID)
		if err != nil {
			return err
		}

		var versionNum int
		if len(latestRevisions) > 0 {
			versionNum = latestRevisions[0].VersionNum + 1
		} else {
			versionNum = 1
		}

		revision := &models.Revision{
			ContentID:   content.ID,
			Title:       content.Title,
			Body:        content.Body,
			Excerpt:     content.Excerpt,
			Status:      models.ContentStatusPublished,
			VersionNum:  versionNum,
			ChangedByID: userID,
		}
		if err := txRepo.CreateRevision(revision); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, ErrContentUpdateFailed
	}

	// Get updated content
	return s.repo.GetContentByID(id)
}

func (s *Service) UnpublishContent(id uint, userID uint) (*models.Content, error) {
	content, err := s.repo.GetContentByID(id)
	if err != nil {
		return nil, err
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := NewRepository(tx)

		// Unpublish content
		if err := txRepo.UnpublishContent(content); err != nil {
			return err
		}

		// Create new revision
		latestRevisions, err := txRepo.GetContentRevisions(content.ID)
		if err != nil {
			return err
		}

		var versionNum int
		if len(latestRevisions) > 0 {
			versionNum = latestRevisions[0].VersionNum + 1
		} else {
			versionNum = 1
		}

		revision := &models.Revision{
			ContentID:   content.ID,
			Title:       content.Title,
			Body:        content.Body,
			Excerpt:     content.Excerpt,
			Status:      models.ContentStatusDraft,
			VersionNum:  versionNum,
			ChangedByID: userID,
		}
		if err := txRepo.CreateRevision(revision); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, ErrContentUpdateFailed
	}

	// Get updated content
	return s.repo.GetContentByID(id)
}

func (s *Service) GetContentRevisions(contentID uint) ([]models.Revision, error) {
	_, err := s.repo.GetContentByID(contentID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetContentRevisions(contentID)
}

func (s *Service) GetContentRevision(contentID uint, revisionID uint) (*models.Revision, error) {
	return s.repo.GetContentRevision(contentID, revisionID)
}

// Category Service Methods
func (s *Service) CreateCategory(req *CreateCategoryRequest) (*models.Category, error) {
	// Check if slug already exists
	existingCategory, err := s.repo.GetCategoryBySlug(req.Slug)
	if err == nil && existingCategory != nil {
		return nil, ErrCategorySlugExists
	}

	category := &models.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
	}

	err = s.repo.CreateCategory(category)
	if err != nil {
		return nil, ErrCategoryCreateFailed
	}

	return category, nil
}

func (s *Service) GetCategory(id uint) (*models.Category, error) {
	return s.repo.GetCategoryByID(id)
}

func (s *Service) UpdateCategory(id uint, req *UpdateCategoryRequest) (*models.Category, error) {
	category, err := s.repo.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}

	// Check if slug is being changed and if it already exists
	if req.Slug != "" && req.Slug != category.Slug {
		existingCategory, err := s.repo.GetCategoryBySlug(req.Slug)
		if err == nil && existingCategory != nil && existingCategory.ID != id {
			return nil, ErrCategorySlugExists
		}
	}

	// Update category fields
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Slug != "" {
		category.Slug = req.Slug
	}
	if req.Description != "" {
		category.Description = req.Description
	}

	err = s.repo.UpdateCategory(category)
	if err != nil {
		return nil, ErrCategoryUpdateFailed
	}

	return category, nil
}

func (s *Service) DeleteCategory(id uint) error {
	_, err := s.repo.GetCategoryByID(id)
	if err != nil {
		return err
	}

	err = s.repo.DeleteCategory(id)
	if err != nil {
		return ErrCategoryDeleteFailed
	}

	return nil
}

func (s *Service) ListCategories() ([]models.Category, error) {
	return s.repo.ListCategories()
}

// Tag Service Methods
func (s *Service) CreateTag(req *CreateTagRequest) (*models.Tag, error) {
	// Check if slug already exists
	existingTag, err := s.repo.GetTagBySlug(req.Slug)
	if err == nil && existingTag != nil {
		return nil, ErrTagSlugExists
	}

	tag := &models.Tag{
		Name: req.Name,
		Slug: req.Slug,
	}

	err = s.repo.CreateTag(tag)
	if err != nil {
		return nil, ErrTagCreateFailed
	}

	return tag, nil
}

func (s *Service) GetTag(id uint) (*models.Tag, error) {
	return s.repo.GetTagByID(id)
}

func (s *Service) UpdateTag(id uint, req *UpdateTagRequest) (*models.Tag, error) {
	tag, err := s.repo.GetTagByID(id)
	if err != nil {
		return nil, err
	}

	// Check if slug is being changed and if it already exists
	if req.Slug != "" && req.Slug != tag.Slug {
		existingTag, err := s.repo.GetTagBySlug(req.Slug)
		if err == nil && existingTag != nil && existingTag.ID != id {
			return nil, ErrTagSlugExists
		}
	}

	// Update tag fields
	if req.Name != "" {
		tag.Name = req.Name
	}
	if req.Slug != "" {
		tag.Slug = req.Slug
	}

	err = s.repo.UpdateTag(tag)
	if err != nil {
		return nil, ErrTagUpdateFailed
	}

	return tag, nil
}

func (s *Service) DeleteTag(id uint) error {
	_, err := s.repo.GetTagByID(id)
	if err != nil {
		return err
	}

	err = s.repo.DeleteTag(id)
	if err != nil {
		return ErrTagDeleteFailed
	}

	return nil
}

func (s *Service) ListTags() ([]models.Tag, error) {
	return s.repo.ListTags()
}

// Media Service Methods
func (s *Service) UploadMedia(file *multipart.FileHeader, req *UploadMediaRequest, userID uint) (*models.Media, error) {
	// Validate file size (10MB limit as an example)
	if file.Size > 10*1024*1024 {
		return nil, ErrFileTooLarge
	}

	// Create uploads directory if it doesn't exist
	uploadsDir := "./uploads/media"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return nil, ErrMediaCreateFailed
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), utils.RandomString(10), ext)
	filePath := filepath.Join(uploadsDir, filename)

	// Open source file
	src, err := file.Open()
	if err != nil {
		return nil, ErrMediaCreateFailed
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, ErrMediaCreateFailed
	}
	defer dst.Close()

	// Copy file content
	if _, err = io.Copy(dst, src); err != nil {
		return nil, ErrMediaCreateFailed
	}

	// Get file type (MIME type)
	fileType := file.Header.Get("Content-Type")
	if fileType == "" {
		fileType = "application/octet-stream"
	}

	// Create media record
	media := &models.Media{
		FileName:       file.Filename,
		FilePath:       "/uploads/media/" + filename,
		FileType:       fileType,
		FileSize:       file.Size,
		MediaType:      req.MediaType,
		AltText:        req.AltText,
		Description:    req.Description,
		UploadedByID:   userID,
		StorageBackend: "local", // Default to local storage for now
	}

	// If media type is image, extract dimensions (would require an image processing library)
	if strings.HasPrefix(fileType, "image/") {
		// Example placeholder: In a real implementation, use an image library to get dimensions
		media.Width = 0
		media.Height = 0
	}

	// If media type is video or audio, extract duration (would require a media processing library)
	if strings.HasPrefix(fileType, "video/") || strings.HasPrefix(fileType, "audio/") {
		// Example placeholder: In a real implementation, use a media library to get duration
		media.Duration = 0
	}

	err = s.repo.CreateMedia(media)
	if err != nil {
		// Clean up file if database save fails
		os.Remove(filePath)
		return nil, ErrMediaCreateFailed
	}

	return media, nil
}

func (s *Service) GetMedia(id uint) (*models.Media, error) {
	return s.repo.GetMediaByID(id)
}

func (s *Service) DeleteMedia(id uint) error {
	media, err := s.repo.GetMediaByID(id)
	if err != nil {
		return err
	}

	// Delete the file
	filePath := "." + media.FilePath
	if err := os.Remove(filePath); err != nil {
		// Log the error but continue to remove the database record
		fmt.Printf("Error deleting media file: %v\n", err)
	}

	// Delete from database
	err = s.repo.DeleteMedia(id)
	if err != nil {
		return ErrMediaDeleteFailed
	}

	return nil
}

func (s *Service) ListMedia(query *MediaListQuery) ([]models.Media, *PaginationResponse, error) {
	// Set default values if not provided
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 100 {
		query.PageSize = 10
	}

	media, total, err := s.repo.ListMedia(query)
	if err != nil {
		return nil, nil, err
	}

	// Calculate pagination info
	totalPages := (int(total) + query.PageSize - 1) / query.PageSize
	pagination := &PaginationResponse{
		Total:       total,
		PerPage:     query.PageSize,
		CurrentPage: query.Page,
		LastPage:    totalPages,
	}

	return media, pagination, nil
}

// Role Service Methods
func (s *Service) CreateRole(req *CreateRoleRequest) (*models.Role, error) {
	// Check if role name already exists
	existingRole, err := s.repo.GetRoleByName(req.Name)
	if err == nil && existingRole != nil {
		return nil, ErrRoleNameExists
	}

	role := &models.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	err = s.repo.CreateRole(role)
	if err != nil {
		return nil, ErrRoleCreateFailed
	}

	return role, nil
}

func (s *Service) GetRole(id uint) (*models.Role, error) {
	return s.repo.GetRoleByID(id)
}

func (s *Service) UpdateRole(id uint, req *UpdateRoleRequest) (*models.Role, error) {
	role, err := s.repo.GetRoleByID(id)
	if err != nil {
		return nil, err
	}

	// Check if name is being changed and if it already exists
	if req.Name != "" && req.Name != role.Name {
		existingRole, err := s.repo.GetRoleByName(req.Name)
		if err == nil && existingRole != nil && existingRole.ID != id {
			return nil, ErrRoleNameExists
		}
	}

	// Update role fields
	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}

	err = s.repo.UpdateRole(role)
	if err != nil {
		return nil, ErrRoleUpdateFailed
	}

	return role, nil
}

func (s *Service) DeleteRole(id uint) error {
	_, err := s.repo.GetRoleByID(id)
	if err != nil {
		return err
	}

	err = s.repo.DeleteRole(id)
	if err != nil {
		return ErrRoleDeleteFailed
	}

	return nil
}

func (s *Service) ListRoles() ([]models.Role, error) {
	return s.repo.ListRoles()
}

func (s *Service) AssignPermissionToRole(roleID uint, permissionID uint) error {
	// Check if role exists
	_, err := s.repo.GetRoleByID(roleID)
	if err != nil {
		return err
	}

	// Check if permission exists
	_, err = s.repo.GetPermissionByID(permissionID)
	if err != nil {
		return err
	}

	err = s.repo.AssignPermissionToRole(roleID, permissionID)
	if err != nil {
		return ErrRolePermissionFailed
	}

	return nil
}

func (s *Service) RemovePermissionFromRole(roleID uint, permissionID uint) error {
	// Check if role exists
	_, err := s.repo.GetRoleByID(roleID)
	if err != nil {
		return err
	}

	// Check if permission exists
	_, err = s.repo.GetPermissionByID(permissionID)
	if err != nil {
		return err
	}

	err = s.repo.RemovePermissionFromRole(roleID, permissionID)
	if err != nil {
		return ErrRolePermissionFailed
	}

	return nil
}

// Permission Service Methods
func (s *Service) CreatePermission(req *CreatePermissionRequest) (*models.Permission, error) {
	// Check if permission name already exists
	existingPermission, err := s.repo.GetPermissionByName(req.Name)
	if err == nil && existingPermission != nil {
		return nil, ErrPermissionNameExists
	}

	permission := &models.Permission{
		Name:        req.Name,
		Description: req.Description,
		Resource:    req.Resource,
		Action:      req.Action,
	}

	err = s.repo.CreatePermission(permission)
	if err != nil {
		return nil, ErrPermissionCreateFailed
	}

	return permission, nil
}

func (s *Service) GetPermission(id uint) (*models.Permission, error) {
	return s.repo.GetPermissionByID(id)
}

func (s *Service) UpdatePermission(id uint, req *UpdatePermissionRequest) (*models.Permission, error) {
	permission, err := s.repo.GetPermissionByID(id)
	if err != nil {
		return nil, err
	}

	// Check if name is being changed and if it already exists
	if req.Name != "" && req.Name != permission.Name {
		existingPermission, err := s.repo.GetPermissionByName(req.Name)
		if err == nil && existingPermission != nil && existingPermission.ID != id {
			return nil, ErrPermissionNameExists
		}
	}

	// Update permission fields
	if req.Name != "" {
		permission.Name = req.Name
	}
	if req.Description != "" {
		permission.Description = req.Description
	}
	if req.Resource != "" {
		permission.Resource = req.Resource
	}
	if req.Action != "" {
		permission.Action = req.Action
	}

	err = s.repo.UpdatePermission(permission)
	if err != nil {
		return nil, ErrPermissionUpdateFailed
	}

	return permission, nil
}

func (s *Service) DeletePermission(id uint) error {
	_, err := s.repo.GetPermissionByID(id)
	if err != nil {
		return err
	}

	err = s.repo.DeletePermission(id)
	if err != nil {
		return ErrPermissionDeleteFailed
	}

	return nil
}

func (s *Service) ListPermissions() ([]models.Permission, error) {
	return s.repo.ListPermissions()
}

// User Role Service Methods
func (s *Service) AssignRoleToUser(userID uint, roleID uint) error {
	// Check if role exists
	_, err := s.repo.GetRoleByID(roleID)
	if err != nil {
		return err
	}

	// Note: In a real implementation, should also check if user exists

	err = s.repo.AssignRoleToUser(userID, roleID)
	if err != nil {
		return ErrUserRoleCreateFailed
	}

	return nil
}

func (s *Service) RemoveRoleFromUser(userID uint, roleID uint) error {
	err := s.repo.RemoveRoleFromUser(userID, roleID)
	if err != nil {
		return ErrUserRoleDeleteFailed
	}

	return nil
}

func (s *Service) GetUserRoles(userID uint) ([]models.Role, error) {
	return s.repo.GetUserRoles(userID)
}

func (s *Service) CheckUserPermission(userID uint, resource string, action string) (bool, error) {
	return s.repo.CheckUserPermission(userID, resource, action)
}
