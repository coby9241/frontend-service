package rbac

import (
	repo "github.com/coby9241/frontend-service/internal/repository/permissions"
	"github.com/qor/roles"
)

// ResourceRBAC defines the RBAC for resource type provided by its string name
func ResourceRBAC(resourceName string, permRepo repo.Repository) (*roles.Permission, error) {
	roleAttrs, err := permRepo.GetPermissionsForResource(resourceName)
	if err != nil {
		return nil, err
	}

	// init roles and deny all
	r := roles.New().NewPermission()

	for _, roleAttr := range roleAttrs {
		if roleAttr.CanCreate {
			r.Allow(roles.Create, roleAttr.RoleName)
		}

		if roleAttr.CanRead {
			r.Allow(roles.Read, roleAttr.RoleName)
		}

		if roleAttr.CanUpdate {
			r.Allow(roles.Update, roleAttr.RoleName)
		}

		if roleAttr.CanDelete {
			r.Allow(roles.Delete, roleAttr.RoleName)
		}
	}

	return r, nil
}
