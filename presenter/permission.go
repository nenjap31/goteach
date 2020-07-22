package presenter

type Permission struct {
	ID     int    `json:"id"`
	Access string `json:"access"`
}

type RolePermission struct {
	ID     int    `json:"id"`
	Access string `json:"access"`
}

func (RolePermission) TableName() string {
	return "role_permission"
}
