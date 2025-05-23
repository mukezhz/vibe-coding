package cms

import (
	"clean-architecture/domain/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Content Repository Methods
func (r *Repository) CreateContent(content *models.Content) error {
	return r.db.Create(content).Error
}

func (r *Repository) GetContentByID(id uint) (*models.Content, error) {
	var content models.Content
	err := r.db.Preload("Categories").Preload("Tags").First(&content, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrContentNotFound
		}
		return nil, err
	}
	return &content, nil
}

func (r *Repository) GetContentBySlug(slug string) (*models.Content, error) {
	var content models.Content
	err := r.db.Preload("Categories").Preload("Tags").Where("slug = ?", slug).First(&content).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrContentNotFound
		}
		return nil, err
	}
	return &content, nil
}

func (r *Repository) UpdateContent(content *models.Content) error {
	return r.db.Save(content).Error
}

func (r *Repository) DeleteContent(id uint) error {
	return r.db.Delete(&models.Content{}, id).Error
}

func (r *Repository) ListContent(query *ContentListQuery) ([]models.Content, int64, error) {
	var contents []models.Content
	var total int64

	db := r.db.Model(&models.Content{})

	// Apply filters
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}
	if query.Search != "" {
		search := "%" + query.Search + "%"
		db = db.Where("title LIKE ? OR body LIKE ? OR excerpt LIKE ?", search, search, search)
	}

	// Count total before pagination
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if query.SortBy != "" {
		direction := "ASC"
		if query.SortDir == "desc" {
			direction = "DESC"
		}
		db = db.Order(query.SortBy + " " + direction)
	} else {
		db = db.Order("created_at DESC")
	}

	// Apply pagination
	offset := (query.Page - 1) * query.PageSize
	err = db.Offset(offset).Limit(query.PageSize).Preload("Categories").Preload("Tags").Find(&contents).Error
	if err != nil {
		return nil, 0, err
	}

	return contents, total, nil
}

func (r *Repository) PublishContent(content *models.Content) error {
	now := time.Now()
	content.Status = models.ContentStatusPublished
	content.PublishedAt = &now
	return r.db.Save(content).Error
}

func (r *Repository) UnpublishContent(content *models.Content) error {
	content.Status = models.ContentStatusDraft
	return r.db.Save(content).Error
}

func (r *Repository) CreateRevision(revision *models.Revision) error {
	return r.db.Create(revision).Error
}

func (r *Repository) GetContentRevisions(contentID uint) ([]models.Revision, error) {
	var revisions []models.Revision
	err := r.db.Where("content_id = ?", contentID).Order("version_num DESC").Find(&revisions).Error
	return revisions, err
}

func (r *Repository) GetContentRevision(contentID uint, revisionID uint) (*models.Revision, error) {
	var revision models.Revision
	err := r.db.Where("content_id = ? AND id = ?", contentID, revisionID).First(&revision).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRevisionNotFound
		}
		return nil, err
	}
	return &revision, nil
}

// Category Repository Methods
func (r *Repository) CreateCategory(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *Repository) GetCategoryByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *Repository) GetCategoryBySlug(slug string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("slug = ?", slug).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *Repository) UpdateCategory(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *Repository) DeleteCategory(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}

func (r *Repository) ListCategories() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

// Tag Repository Methods
func (r *Repository) CreateTag(tag *models.Tag) error {
	return r.db.Create(tag).Error
}

func (r *Repository) GetTagByID(id uint) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.First(&tag, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTagNotFound
		}
		return nil, err
	}
	return &tag, nil
}

func (r *Repository) GetTagBySlug(slug string) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.Where("slug = ?", slug).First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTagNotFound
		}
		return nil, err
	}
	return &tag, nil
}

func (r *Repository) UpdateTag(tag *models.Tag) error {
	return r.db.Save(tag).Error
}

func (r *Repository) DeleteTag(id uint) error {
	return r.db.Delete(&models.Tag{}, id).Error
}

func (r *Repository) ListTags() ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.Find(&tags).Error
	return tags, err
}

// Media Repository Methods
func (r *Repository) CreateMedia(media *models.Media) error {
	return r.db.Create(media).Error
}

func (r *Repository) GetMediaByID(id uint) (*models.Media, error) {
	var media models.Media
	err := r.db.First(&media, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMediaNotFound
		}
		return nil, err
	}
	return &media, nil
}

func (r *Repository) DeleteMedia(id uint) error {
	return r.db.Delete(&models.Media{}, id).Error
}

func (r *Repository) ListMedia(query *MediaListQuery) ([]models.Media, int64, error) {
	var media []models.Media
	var total int64

	db := r.db.Model(&models.Media{})

	// Apply filters
	if query.MediaType != "" {
		db = db.Where("media_type = ?", query.MediaType)
	}
	if query.Search != "" {
		search := "%" + query.Search + "%"
		db = db.Where("file_name LIKE ? OR description LIKE ? OR alt_text LIKE ?", search, search, search)
	}

	// Count total before pagination
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if query.SortBy != "" {
		direction := "ASC"
		if query.SortDir == "desc" {
			direction = "DESC"
		}
		db = db.Order(query.SortBy + " " + direction)
	} else {
		db = db.Order("created_at DESC")
	}

	// Apply pagination
	offset := (query.Page - 1) * query.PageSize
	err = db.Offset(offset).Limit(query.PageSize).Find(&media).Error
	if err != nil {
		return nil, 0, err
	}

	return media, total, nil
}

// Role Repository Methods
func (r *Repository) CreateRole(role *models.Role) error {
	return r.db.Create(role).Error
}

func (r *Repository) GetRoleByID(id uint) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("Permissions").First(&role, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}
	return &role, nil
}

func (r *Repository) GetRoleByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}
	return &role, nil
}

func (r *Repository) UpdateRole(role *models.Role) error {
	return r.db.Save(role).Error
}

func (r *Repository) DeleteRole(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}

func (r *Repository) ListRoles() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Preload("Permissions").Find(&roles).Error
	return roles, err
}

func (r *Repository) AssignPermissionToRole(roleID, permissionID uint) error {
	return r.db.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", roleID, permissionID).Error
}

func (r *Repository) RemovePermissionFromRole(roleID, permissionID uint) error {
	return r.db.Exec("DELETE FROM role_permissions WHERE role_id = ? AND permission_id = ?", roleID, permissionID).Error
}

// Permission Repository Methods
func (r *Repository) CreatePermission(permission *models.Permission) error {
	return r.db.Create(permission).Error
}

func (r *Repository) GetPermissionByID(id uint) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.First(&permission, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPermissionNotFound
		}
		return nil, err
	}
	return &permission, nil
}

func (r *Repository) GetPermissionByName(name string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("name = ?", name).First(&permission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPermissionNotFound
		}
		return nil, err
	}
	return &permission, nil
}

func (r *Repository) UpdatePermission(permission *models.Permission) error {
	return r.db.Save(permission).Error
}

func (r *Repository) DeletePermission(id uint) error {
	return r.db.Delete(&models.Permission{}, id).Error
}

func (r *Repository) ListPermissions() ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Find(&permissions).Error
	return permissions, err
}

// User Role Repository Methods
func (r *Repository) AssignRoleToUser(userID, roleID uint) error {
	userRole := models.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	return r.db.Create(&userRole).Error
}

func (r *Repository) RemoveRoleFromUser(userID, roleID uint) error {
	return r.db.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&models.UserRole{}).Error
}

func (r *Repository) GetUserRoles(userID uint) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Preload("Permissions").
		Find(&roles).Error
	return roles, err
}

func (r *Repository) CheckUserPermission(userID uint, resource string, action string) (bool, error) {
	var count int64
	err := r.db.Raw(`
		SELECT COUNT(*) FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = ? AND p.resource = ? AND p.action = ?
	`, userID, resource, action).Count(&count).Error

	return count > 0, err
}
