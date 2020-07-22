package presenter

import "goteach/models"

type User struct {
	models.BaseModel
	Name        string       `json:"name"`
	Username    string       `json:"username"`
	Email       string       `json:"email"`
	Role        models.Role  `json:"role"`
	RoleID      int          `json:"role_id"`
	IsActive    bool         `json:"is_active"`
	Permissions []Permission `json:"permissions"`
}
