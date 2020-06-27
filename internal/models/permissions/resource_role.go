package permissions

// ResourceRole is
type ResourceRole struct {
	CanCreate bool `gorm:"not null,DEFAULT false"`
	CanRead   bool `gorm:"not null,DEFAULT false"`
	CanUpdate bool `gorm:"not null,DEFAULT false"`
	CanDelete bool `gorm:"not null,DEFAULT false"`
}

// TableName sets ResourceRole's table name
func (ResourceRole) TableName() string {
	return "resource_role"
}
