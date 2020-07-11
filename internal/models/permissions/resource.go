package permissions

import "time"

// IResource is an interface to specify a particular model is also a resource
type IResource interface {
	GetResourceName() string
}

// Resource is an struct that holds the definition of a resource that links to the RBAC
type Resource struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string  `sql:"unique_index:uix_resource_name"`
	Roles     []*Role `gorm:"many2many:resource_role;association_foreignkey:ID;foreignkey:ID"`
}
