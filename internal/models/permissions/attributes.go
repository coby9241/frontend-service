package permissions

// EnabledAttributes is
type EnabledAttributes struct {
	CanCreate bool
	CanRead   bool
	CanUpdate bool
	CanDelete bool
}

// RoleAttributes is
type RoleAttributes struct {
	RoleName  string
	CanCreate bool
	CanRead   bool
	CanUpdate bool
	CanDelete bool
}
