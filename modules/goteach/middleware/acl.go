package middleware

import (
	"github.com/davecgh/go-spew/spew"
	"goteach/config"
	"goteach/presenter"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

const (
	UPDATE_JOB = "update_job"
)

// ACL is method for checking user permisson
func ACL(permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if isAdmin(c) {
				return next(c)
			}
			db := config.DB
			var userPermissions []presenter.RolePermission
			roleID := getRoleID(c)
			db = db.Where("role_id = ?", roleID)
			db = db.Select([]string{"permission_id as id", "concat(p.permission, '_', p.resource) as access"})
			db.Joins("left join permissions p on role_permission.permission_id = p.id").Find(&userPermissions)
			spew.Dump(userPermissions)
			for _, p := range userPermissions {
				if permission == p.Access {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, "You don't have permission to access this resource")
		}
	}
}

func isAdmin(c echo.Context) bool {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*config.JwtCustomClaims)
	return claims.IsAdmin
}

func getRoleID(c echo.Context) int {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*config.JwtCustomClaims)
	return claims.RoleID
}
