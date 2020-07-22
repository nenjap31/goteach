package models

import "time"

type Permission struct {
	ID          int       `json:"id" gorm:"primary_key"`
	Resource    string    `json:"resource" gorm:"not null; size:50"`
	Permission  string    `json:"permission" gorm:"not null; size:50"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
