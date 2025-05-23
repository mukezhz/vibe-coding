package cms

import (
	"clean-architecture/domain/models"
	"clean-architecture/pkg/framework"

	"gorm.io/gorm"
)

// Migrate creates the necessary tables for the CMS module
func Migrate(db *gorm.DB, logger framework.Logger) {
	logger.Info("Running CMS migrations")

	// Create tables if they don't exist
	err := db.AutoMigrate(
		&models.Content{},
		&models.Category{},
		&models.Tag{},
		&models.Revision{},
		&models.Media{},
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
	)
	if err != nil {
		logger.Error("Failed to migrate CMS tables: " + err.Error())
		return
	}

	// Create default roles and permissions
	err = CreateDefaultRolesAndPermissions(db)
	if err != nil {
		logger.Error("Failed to create default roles and permissions: " + err.Error())
		return
	}

	logger.Info("CMS migrations completed successfully")
}

// CreateDefaultRolesAndPermissions creates default roles and permissions for the CMS
func CreateDefaultRolesAndPermissions(db *gorm.DB) error {
	// Create default roles
	adminRole := models.Role{
		Name:        "admin",
		Description: "Administrator with full access",
	}

	editorRole := models.Role{
		Name:        "editor",
		Description: "Can create and edit content",
	}

	authorRole := models.Role{
		Name:        "author",
		Description: "Can create content but not publish",
	}

	// Check if admin role already exists
	var existingAdminRole models.Role
	result := db.Where("name = ?", "admin").First(&existingAdminRole)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	if result.Error == gorm.ErrRecordNotFound {
		// Create admin role
		if err := db.Create(&adminRole).Error; err != nil {
			return err
		}
	} else {
		adminRole = existingAdminRole
	}

	// Check if editor role already exists
	var existingEditorRole models.Role
	result = db.Where("name = ?", "editor").First(&existingEditorRole)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	if result.Error == gorm.ErrRecordNotFound {
		// Create editor role
		if err := db.Create(&editorRole).Error; err != nil {
			return err
		}
	} else {
		editorRole = existingEditorRole
	}

	// Check if author role already exists
	var existingAuthorRole models.Role
	result = db.Where("name = ?", "author").First(&existingAuthorRole)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	if result.Error == gorm.ErrRecordNotFound {
		// Create author role
		if err := db.Create(&authorRole).Error; err != nil {
			return err
		}
	} else {
		authorRole = existingAuthorRole
	}

	// Define permissions
	permissions := []models.Permission{
		{Name: "create_content", Description: "Can create content", Resource: "content", Action: "create"},
		{Name: "read_content", Description: "Can read content", Resource: "content", Action: "read"},
		{Name: "update_content", Description: "Can update content", Resource: "content", Action: "update"},
		{Name: "delete_content", Description: "Can delete content", Resource: "content", Action: "delete"},
		{Name: "publish_content", Description: "Can publish content", Resource: "content", Action: "publish"},
		{Name: "create_category", Description: "Can create categories", Resource: "category", Action: "create"},
		{Name: "read_category", Description: "Can read categories", Resource: "category", Action: "read"},
		{Name: "update_category", Description: "Can update categories", Resource: "category", Action: "update"},
		{Name: "delete_category", Description: "Can delete categories", Resource: "category", Action: "delete"},
		{Name: "create_tag", Description: "Can create tags", Resource: "tag", Action: "create"},
		{Name: "read_tag", Description: "Can read tags", Resource: "tag", Action: "read"},
		{Name: "update_tag", Description: "Can update tags", Resource: "tag", Action: "update"},
		{Name: "delete_tag", Description: "Can delete tags", Resource: "tag", Action: "delete"},
		{Name: "create_media", Description: "Can upload media", Resource: "media", Action: "create"},
		{Name: "read_media", Description: "Can read media", Resource: "media", Action: "read"},
		{Name: "delete_media", Description: "Can delete media", Resource: "media", Action: "delete"},
		{Name: "create_role", Description: "Can create roles", Resource: "role", Action: "create"},
		{Name: "read_role", Description: "Can read roles", Resource: "role", Action: "read"},
		{Name: "update_role", Description: "Can update roles", Resource: "role", Action: "update"},
		{Name: "delete_role", Description: "Can delete roles", Resource: "role", Action: "delete"},
		{Name: "create_permission", Description: "Can create permissions", Resource: "permission", Action: "create"},
		{Name: "read_permission", Description: "Can read permissions", Resource: "permission", Action: "read"},
		{Name: "update_permission", Description: "Can update permissions", Resource: "permission", Action: "update"},
		{Name: "delete_permission", Description: "Can delete permissions", Resource: "permission", Action: "delete"},
		{Name: "manage_user_roles", Description: "Can manage user roles", Resource: "user_role", Action: "manage"},
	}

	// Create permissions
	for _, permission := range permissions {
		var existingPermission models.Permission
		result := db.Where("name = ?", permission.Name).First(&existingPermission)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			return result.Error
		}

		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&permission).Error; err != nil {
				return err
			}
		}
	}

	// Assign all permissions to admin role
	var allPermissions []models.Permission
	if err := db.Find(&allPermissions).Error; err != nil {
		return err
	}

	if err := db.Model(&adminRole).Association("Permissions").Clear(); err != nil {
		return err
	}
	if err := db.Model(&adminRole).Association("Permissions").Append(allPermissions); err != nil {
		return err
	}

	// Assign content and media permissions to editor role
	var editorPermissions []models.Permission
	if err := db.Where("resource IN ?", []string{"content", "category", "tag", "media"}).Find(&editorPermissions).Error; err != nil {
		return err
	}

	if err := db.Model(&editorRole).Association("Permissions").Clear(); err != nil {
		return err
	}
	if err := db.Model(&editorRole).Association("Permissions").Append(editorPermissions); err != nil {
		return err
	}

	// Assign limited content permissions to author role
	var authorPermissions []models.Permission
	if err := db.Where("name IN ?", []string{"create_content", "read_content", "update_content", "read_category", "read_tag", "create_media", "read_media"}).Find(&authorPermissions).Error; err != nil {
		return err
	}

	if err := db.Model(&authorRole).Association("Permissions").Clear(); err != nil {
		return err
	}
	if err := db.Model(&authorRole).Association("Permissions").Append(authorPermissions); err != nil {
		return err
	}

	return nil
}
