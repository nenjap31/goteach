package routes

import (
	c "goteach/controllers"
	"goteach/config"
	goteachMiddleware "goteach/modules/goteach/middleware"
	
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func acl(permission string) echo.MiddlewareFunc {
	return goteachMiddleware.ACL(permission)
}

func Api(e *echo.Echo) {
	e.POST("login", c.LoginUser)
	e.GET("profile", c.GetProfile, middleware.JWTWithConfig(config.JwtConfig))

	// user
	user := e.Group("user", middleware.JWTWithConfig(config.JwtConfig))
	user.GET("", c.GetUser, acl("read_user"))
	user.POST("", c.AddUser, acl("create_user"))
	user.PUT("/:id", c.UpdateUser, acl("update_user"))
	
}
