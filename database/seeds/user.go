package seeds

import (
	"goteach/models"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	SeederList["User"] = User{}
}

type User struct{}

func (seeder User) Run(db *gorm.DB) error {


	var role models.Role
	var count int
	db.Where("name = ?", "admin").Find(&role).Count(&count)

	if count < 1 {
		RoleSeeder{}.Run(db)
		db.Where("name = ?", "admin").Find(&role)
	}

	var user = models.User{
		Username: "goteach",
		Email:    "goteach@gmail.com",
		Password: "goteach",
		RoleID:   role.ID,
	}
	user.Password, _ = seeder.HashPassword(user.Password)
	db.Save(&user)
	return nil
}

func (User) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(bytes), err
}
