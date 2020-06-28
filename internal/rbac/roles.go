package rbac

import (
	"net/http"

	"github.com/coby9241/frontend-service/internal/models/users"
	repo "github.com/coby9241/frontend-service/internal/repository/permissions"
	"github.com/qor/roles"
)

// Load Register roles on startup
func Load(permRepo repo.Repository) error {
	appRoles, err := permRepo.GetRoles()
	if err != nil {
		return err
	}

	for _, role := range appRoles {
		// prevent shadowing
		r := role
		roles.Register(r.Name, func(req *http.Request, currentUser interface{}) bool {
			usr, ok := currentUser.(*users.User)
			if !ok {
				return false
			}
			return usr.Role.Name == r.Name
		})
	}

	return nil
}
