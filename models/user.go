package models

type User struct {
	BaseModel `structs:",flatten"`
	Name      string `json:"name" gorm:"size:50"`
	Username  string `json:"username" gorm:"not null; unique" gorm:"size:50"`
	Password  string `json:"password" gorm:"not null" gorm:"size:100"`
	Email     string `json:"email" gorm:"not null; unique" gorm:"size:50"`
	RoleID    int    `json:"role_id" gorm:"not null"`
	IsActive  bool   `json:"is_active" sql:"default:1"`
	Role      Role   `json:"role" structs:"-"`
}

