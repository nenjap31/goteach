package controllers

import (
	"encoding/json"
	"flag"
	"goteach/config"
	"goteach/models"
	"goteach/presenter"
	"os"
	"strings"

	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"github.com/thedevsaddam/govalidator"
	"golang.org/x/crypto/bcrypt"
)

const (
	MESSAGE         = "message"
	VALIDATIONERROR = "validationError"
	INVALIDUSER = "Invalid username or password"
	REQUIRED ="required"
	USERNAME = "username"
	ROLE_PERMISSION_TEXT = "Role.Permissions"
)

func LoginUser(c echo.Context) error {
	type LoginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var user models.User
	db := config.DB
	login := new(LoginData)

	if flag.Lookup("test.v") == nil && !strings.HasSuffix(os.Args[0], ".test") {
		_ = c.Bind(login)
	} else {
		_ = json.NewDecoder(c.Request().Body).Decode(&login)
	}
	if db.Preload("Role").Where("username = ?", login.Username).First(&user).RecordNotFound() {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{MESSAGE: INVALIDUSER})
	} else if !user.IsActive {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{MESSAGE: "User disabled"})
	} else if CheckPasswordHash(login.Password, user.Password) {

		// Set custom claims
		claims := &config.JwtCustomClaims{
			user.ID,
			user.Username,
			user.RoleID,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			},
			user.Role.IsAdmin,
		}

		// Create token with claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(viper.GetString("jwtSign")))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, echo.Map{
			"token": t,
		})
	} else {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{MESSAGE: INVALIDUSER})
	}


	return echo.ErrUnauthorized
}


func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetProfile(c echo.Context) error {
	var user presenter.User
	// get database instance
	db := config.DB
	// get user id from token claim
	claim := getUser(c)

	// get user data by user id
	db.First(&user, claim.ID)
	db.Preload(ROLE_PERMISSION_TEXT).Find(&user)


	return c.JSON(http.StatusOK, user)
}

/*func getToken(c echo.Context) (token string) {
	user := c.Get("user").(*jwt.Token)

	return user.Raw
}*/

// getUser() get authecticated user
func getUser(c echo.Context) *config.JwtCustomClaims {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*config.JwtCustomClaims)
	return claims
}

func GetListUser(c echo.Context) error {
	db := config.DB
	var user []models.User

	db.Find(&user)
	return c.JSON(http.StatusOK, user)
}


func AddUser(c echo.Context) error {
	var user models.User
	db := config.DB

	rules := govalidator.MapData{
		USERNAME: []string{REQUIRED},
		"email":    []string{REQUIRED, "email"},
		"password": []string{REQUIRED},
	}

	opts := govalidator.Options{
		Request: c.Request(),
		Data:    &user,
		Rules:   rules,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()

	if len(e) > 0 {
		err := map[string]interface{}{VALIDATIONERROR: e}
		return c.JSON(http.StatusBadRequest, err)
	}

	var users []models.User
	db.Where("username = ? ", user.Username).Find(&users)
	if len(users) > 0 {
		er := map[string]interface{}{MESSAGE: []string{"username already used"}}
		err := map[string]interface{}{VALIDATIONERROR: er}
		return c.JSON(http.StatusBadRequest, err)
	}

	db.Where("email = ? ", user.Email).Find(&users)
	if len(users) > 0 {
		er := map[string]interface{}{MESSAGE: []string{"email already used"}}
		err := map[string]interface{}{VALIDATIONERROR: er}
		return c.JSON(http.StatusBadRequest, err)
	}

	user.Password, _ = HashPassword(user.Password)

	if err := db.Create(&user).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	db.Preload(ROLE_PERMISSION_TEXT).Find(&user)

	return c.JSON(http.StatusCreated, user)
}

func UpdateUser(c echo.Context) error {
	var newData models.User
	db := config.DB

	rules := govalidator.MapData{
		USERNAME: []string{},
		"email":    []string{"email"},
	}

	opts := govalidator.Options{
		Request: c.Request(),
		Data:    &newData,
		Rules:   rules,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()

	if len(e) > 0 {
		err := map[string]interface{}{VALIDATIONERROR: e}
		return c.JSON(http.StatusBadRequest, err)
	}

	var user models.User
	var count int

	id := c.Param("id")
	db.First(&user, id).Count(&count)
	if count != 1 {
		er := map[string]interface{}{MESSAGE: []string{"user does not exist"}}
		err := map[string]interface{}{VALIDATIONERROR: er}
		return c.JSON(http.StatusBadRequest, err)
	}

	var users []models.User
	db.Where("username = ? ", newData.Username).Find(&users)
	if len(users) > 0 {
		er := map[string]interface{}{MESSAGE: []string{"username already used"}}
		err := map[string]interface{}{VALIDATIONERROR: er}
		return c.JSON(http.StatusBadRequest, err)
	}

	db.Where("email = ? ", newData.Email).Find(&users)
	if len(users) > 0 {
		er := map[string]interface{}{MESSAGE: []string{"email already used"}}
		err := map[string]interface{}{VALIDATIONERROR: er}
		return c.JSON(http.StatusBadRequest, err)
	}

	if len(newData.Password) > 0 {
		newData.Password, _ = HashPassword(newData.Password)
	}
	db.Model(&user).Updates(newData)
	db.Model(&user).Update("is_active", newData.IsActive)
	db.Preload(ROLE_PERMISSION_TEXT).Find(&user)
	return c.JSON(http.StatusOK, user)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(bytes), err
}
