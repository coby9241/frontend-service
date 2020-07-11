package migrations

import (
	"time"

	"github.com/coby9241/frontend-service/internal/config"
	"github.com/coby9241/frontend-service/internal/models/permissions"
	"github.com/coby9241/frontend-service/internal/models/users"
	repo "github.com/coby9241/frontend-service/internal/repository/permissions"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// InitSchema is the function to initialise the schema for newly initiated DBs
func InitSchema(tx *gorm.DB) error {
	transaction := tx.Begin()
	err := transaction.AutoMigrate(
		&permissions.Role{},
		&permissions.Resource{},
		&permissions.ResourceRole{},
		&users.User{},
	).Error

	if err != nil {
		return rollbackAndErr(transaction, err)
	}

	// add resource
	resource := permissions.Resource{ResourceName: users.User{}.GetResourceName()}
	if err := transaction.Save(&resource).Error; err != nil {
		return rollbackAndErr(transaction, err)
	}

	// add basic roles
	rolesList := []struct {
		roleName     string
		roleResource []*permissions.Resource
		rolePerms    permissions.EnabledAttributes
	}{
		{
			roleName:     "admin",
			roleResource: []*permissions.Resource{&resource},
			rolePerms: permissions.EnabledAttributes{
				CanCreate: true,
				CanRead:   true,
				CanUpdate: true,
				CanDelete: true,
			},
		},
		{
			roleName:     "editor",
			roleResource: []*permissions.Resource{&resource},
			rolePerms: permissions.EnabledAttributes{
				CanCreate: false,
				CanRead:   true,
				CanUpdate: false,
				CanDelete: false,
			},
		},
		{
			roleName:     "viewer",
			roleResource: []*permissions.Resource{&resource},
			rolePerms: permissions.EnabledAttributes{
				CanCreate: false,
				CanRead:   false,
				CanUpdate: false,
				CanDelete: false,
			},
		},
	}

	permRepo := repo.NewPermissionsRepositoryImpl(transaction)
	// create role from list
	for _, role := range rolesList {
		createdRole, err := permRepo.CreateNewRole(role.roleResource, role.roleName)
		if err != nil {
			return rollbackAndErr(transaction, err)
		}

		if err := permRepo.SetPermissions(resource.ID, createdRole.ID, role.rolePerms); err != nil {
			return rollbackAndErr(transaction, err)
		}
	}

	// add FK resource_role
	transaction.Table("resource_role").AddForeignKey("role_id", "roles(id)", "CASCADE", "CASCADE")
	transaction.Table("resource_role").AddForeignKey("resource_id", "resources(id)", "CASCADE", "CASCADE")

	// add admin user
	var pwd []byte
	if pwd, err = bcrypt.GenerateFromPassword([]byte(config.GetInstance().AdminPassword), bcrypt.DefaultCost); err != nil {
		return rollbackAndErr(transaction, err)
	}

	var adminRole permissions.Role
	if err := transaction.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return rollbackAndErr(transaction, err)
	}

	currTime := time.Now()
	usr := users.User{
		Provider:          "email",
		UID:               config.GetInstance().AdminUsername,
		PasswordHash:      string(pwd),
		UserID:            "admin",
		PasswordChangedAt: &currTime,
		Model: gorm.Model{
			CreatedAt: currTime,
			UpdatedAt: currTime,
		},
		RoleID: adminRole.ID,
	}

	if err = transaction.Create(&usr).Error; err != nil {
		return rollbackAndErr(transaction, err)
	}

	return transaction.Commit().Error
}
