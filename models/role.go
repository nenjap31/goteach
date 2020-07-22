package models

import "time"

type Role struct {
	ID          int          `json:"id" gorm:"primary_key"`
	Name        string       `json:"name" gorm:"not null; unique; size:50"`
	CreatedAt   time.Time    `json:"created_at"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permission"`
	IsAdmin     bool         `json:"is_admin" sql:"default:0"`
}
