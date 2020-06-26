package users

import (
	"github.com/jinzhu/gorm"
)

// Role contains the different RBAC roles for a user
type Role struct {
	gorm.Model
	Name string
}

// TableName Set Role's table name
func (Role) TableName() string {
	return "roles"
}
