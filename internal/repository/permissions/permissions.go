package permissions

import (
	"github.com/coby9241/frontend-service/internal/models/permissions"
	"github.com/jinzhu/gorm"
)

// Repository is
type Repository interface {
	// Create
	CreateNewRole(resources []*permissions.Resource, roleName string) (*permissions.Role, error)
	// Read
	GetRoles() ([]permissions.Role, error)
	GetPermissionsForResource(resourceName string) ([]permissions.RoleAttributes, error)
	// Update
	SetPermissions(resourceID, roleID uint, perm permissions.EnabledAttributes) error
}

// RepositoryImpl is
type RepositoryImpl struct {
	DB *gorm.DB
}

var _ Repository = new(RepositoryImpl)

// NewPermissionsRepositoryImpl is
func NewPermissionsRepositoryImpl(storage *gorm.DB) Repository {
	return &RepositoryImpl{
		DB: storage,
	}
}

// CreateNewRole is
func (p *RepositoryImpl) CreateNewRole(resources []*permissions.Resource, roleName string) (*permissions.Role, error) {
	role := permissions.Role{
		Name:      roleName,
		Resources: resources,
	}

	if err := p.DB.Save(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

// SetPermissions is
func (p *RepositoryImpl) SetPermissions(resourceID, roleID uint, perm permissions.EnabledAttributes) (err error) {
	rr := permissions.ResourceRole{}
	err = p.DB.
		Model(&rr).
		Where("resource_id = ? AND role_id = ?", resourceID, roleID).
		Updates(map[string]interface{}{
			"CanCreate": perm.CanCreate,
			"CanRead":   perm.CanRead,
			"CanUpdate": perm.CanUpdate,
			"CanDelete": perm.CanDelete,
		}).Error

	return
}

// GetRoles is
func (p *RepositoryImpl) GetRoles() ([]permissions.Role, error) {
	roles := []permissions.Role{}
	if err := p.DB.Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

// GetPermissionsForResource is
func (p *RepositoryImpl) GetPermissionsForResource(resourceName string) ([]permissions.RoleAttributes, error) {
	roleAttrs := make([]permissions.RoleAttributes, 0)

	// get all roles for resource
	r := permissions.Resource{}
	var roles []permissions.Role
	if err := p.DB.Where("name = ?", resourceName).First(&r).Error; err != nil {
		return nil, err
	}

	if err := p.DB.Model(&r).Related(&roles, "Roles").Error; err != nil {
		return nil, err
	}

	// get all permissions for role
	for _, role := range roles {
		rr := permissions.ResourceRole{}
		if err := p.DB.Where("resource_id = ? AND role_id = ?", r.ID, role.ID).First(&rr).Error; err != nil {
			return nil, err
		}

		roleAttr := permissions.RoleAttributes{
			RoleName:  role.Name,
			CanCreate: rr.CanCreate,
			CanRead:   rr.CanRead,
			CanUpdate: rr.CanUpdate,
			CanDelete: rr.CanDelete,
		}

		roleAttrs = append(roleAttrs, roleAttr)
	}

	return roleAttrs, nil
}
