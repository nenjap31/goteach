package models

type RolePermission struct {
	RoleID       int
	PermissionID int
}

func (RolePermission) TableName() string {
	return "role_permission"
}
