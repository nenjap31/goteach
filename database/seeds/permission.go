package seeds

import (
	"goteach/models"

	"github.com/jinzhu/gorm"
)

func init() {
	SeederList["PermissionSeeder"] = PermissionSeeder{}
}

type PermissionSeeder struct {
	jsonParser
}

func (seeder PermissionSeeder) Run(db *gorm.DB) error {
	var permissions []models.Permission
	seeder.ParseJSON("permissions.json", &permissions)
	for _, permission := range permissions {
		db.FirstOrCreate(&permission, permission)
	}
	return nil
}
