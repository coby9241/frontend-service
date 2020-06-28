package permissions

import "time"

// Role contains the different RBAC roles for a user
type Role struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string      `sql:"type:varchar;unique_index:uix_role_name"`
	Resources []*Resource `gorm:"many2many:resource_role;association_foreignkey:ID;foreignkey:ID"`
}

// TableName Set Role's table name
func (Role) TableName() string {
	return "roles"
}
