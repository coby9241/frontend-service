package permissions

import (
	"github.com/jinzhu/gorm"
)

// Role contains the different RBAC roles for a user
type Role struct {
	gorm.Model
	Name      string      `sql:"type:varchar;unique_index:uix_role_name"`
	Resources []*Resource `gorm:"many2many:resource_role"`
}

// TableName Set Role's table name
func (Role) TableName() string {
	return "roles"
}
