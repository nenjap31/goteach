package seeds

import (
	"goteach/models"

	"github.com/jinzhu/gorm"
)

func init() {
	SeederList["RoleSeeder"] = RoleSeeder{}
}

type RoleSeeder struct{}

func (RoleSeeder) Run(db *gorm.DB) error {
	var role models.Role
	role.Name = "admin"
	role.IsAdmin = true

	if err := db.Create(&role).Error; err != nil {
		return err
	}
	return nil
}
