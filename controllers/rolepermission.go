package controllers

import (
	"encoding/json"
	"flag"
	"goteach/config"
	"goteach/models"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo"
)

type RoleData struct {
	Name        string `json:"name"`
	Isadmin bool `json:"is_admin"`
	Permissions []int  `json:"permissions"`
}

const (
	PERMISSION_TEXT ="Permissions"
)

func GetRole(c echo.Context) error {
	var roles []models.Role
	db := config.DB

	db.Preload(PERMISSION_TEXT).Find(&roles)

	return c.JSON(http.StatusOK, roles)
}

func GetPermission(c echo.Context) error {
	db := config.DB
	var permissions []models.Permission

	db.Find(&permissions)
	return c.JSON(http.StatusOK, permissions)
}

func AddRole(c echo.Context) error {
	db := config.DB
	tx := db.Begin()
	var postData RoleData
	var role models.Role

	if flag.Lookup("test.v") == nil && !strings.HasSuffix(os.Args[0], ".test") {
		_ = c.Bind(&postData)
	} else {
		_ = json.NewDecoder(c.Request().Body).Decode(&postData)
	}
	role.Name = postData.Name
	var roles []models.Role
	db.Where("name = ? ", postData.Name).Find(&roles)
	if len(roles) > 0 {
		er := map[string]interface{}{MESSAGE: []string{"role name already used"}}
		err := map[string]interface{}{VALIDATIONERROR: er}
		return c.JSON(http.StatusBadRequest, err)
	}

	tx.Save(&role)
	for _, pid := range postData.Permissions {
		tx.Create(models.RolePermission{RoleID: role.ID, PermissionID: pid})
	}
	tx.Commit()

	db.Preload(PERMISSION_TEXT).Find(&role)

	return c.JSON(http.StatusCreated, role)
}

func UpdateRole(c echo.Context) error {
	db := config.DB
	var role models.Role
	var newData RoleData
	id := c.Param("id")

	if flag.Lookup("test.v") == nil && !strings.HasSuffix(os.Args[0], ".test") {
		_ = c.Bind(&newData)
	} else {
		_ = json.NewDecoder(c.Request().Body).Decode(&newData)
	}
	if db.First(&role, id).RecordNotFound() {
		er := map[string]interface{}{MESSAGE: []string{"invalid role id"}}
		err := map[string]interface{}{VALIDATIONERROR: er}
		return c.JSON(http.StatusBadRequest, err)
	}

	var roles []models.Role
	db.Where("name = ? ", newData.Name).Where("id != ? ", role.ID).Find(&roles)
	if len(roles) > 0 {

		er := map[string]interface{}{MESSAGE: []string{"role name already used"}}
		err := map[string]interface{}{VALIDATIONERROR: er}
		return c.JSON(http.StatusBadRequest, err)
	}

	if newData.Name != "" {
		db.Model(&role).Updates(map[string]interface{}{"name": newData.Name,"is_admin":newData.Isadmin})
	}

	db.Delete(models.RolePermission{}, "role_id = ?", role.ID)
	for _, pid := range newData.Permissions {
		db.Create(models.RolePermission{RoleID: role.ID, PermissionID: pid})
	}

	db.Preload(PERMISSION_TEXT).Find(&role)

	return c.JSON(http.StatusOK, role)
}
