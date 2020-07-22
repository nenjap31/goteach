package controllers

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	mocket "github.com/selvatico/go-mocket"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"goteach/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRolePermission(t *testing.T){
	t.Run("test.v", func(t *testing.T) {
		e := echo.New()

		claims := &config.JwtCustomClaims{
			ID:       1,
			Username: "admin",
			RoleID:   1,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			},
			IsAdmin: true,
		}
		encToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signToken, _ := encToken.SignedString([]byte(viper.GetString("jwtSign")))

		//test get role list
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+ signToken)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("role")

		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  `SELECT * FROM "role"  WHERE`, // the same as .WithQuery()
				Response: Roleresp,                 // the same as .WithReply
				Once:     true,                             // To not use it twice if true
			},

		})

		if assert.NoError(t, GetRole(c)) {
			spew.Dump("getrole========")
			spew.Dump(rec.Code)
			assert.Equal(t, http.StatusOK, rec.Code)
		}


		//test get permission list
		req = httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+ signToken)
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.SetPath("permission")

		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  `SELECT * FROM "role"  WHERE`, // the same as .WithQuery()
				Response: Roleresp,                 // the same as .WithReply
				Once:     true,                             // To not use it twice if true
			},

		})

		if assert.NoError(t, GetPermission(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}

		//test add role success
		req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(roleJSON))
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+ signToken)
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.SetPath("role")

		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  `SELECT * FROM "roles"  WHERE`, // the same as .WithQuery()
				Response: Rolerespfail,                 // the same as .WithReply
				Once:     true,                             // To not use it twice if true
			},
			{
				Pattern:  `INSERT INTO "role_permission"`, // the same as .WithQuery()
				Response: Roleresp,                 // the same as .WithReply
				Once:     true,                             // To not use it twice if true
			},

		})

		if assert.NoError(t, AddRole(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}


		//test add role failed, role name already exist
		req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(roleJSON))
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+ signToken)
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.SetPath("role")

		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  `SELECT * FROM "roles"  WHERE`, // the same as .WithQuery()
				Response: Rolerespfail2,                 // the same as .WithReply
				Once:     true,                             // To not use it twice if true
			},
			{
				Pattern:  `INSERT INTO "role_permission"`, // the same as .WithQuery()
				Response: Roleresp,                 // the same as .WithReply
				Once:     true,                             // To not use it twice if true
			},

		})

		if assert.NoError(t, AddRole(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}

		//test update role success
		req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(roleJSON))
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+ signToken)
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.SetPath("role/1")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  `SELECT * FROM "roles"  WHERE`, // the same as .WithQuery()
				Response: Roleresp,                 // the same as .WithReply
				Once:     true,                             // To not use it twice if true
			},

			{
				Pattern:  `UPDATE "roles" SET`, // the same as .WithQuery()
				Response: Roleresp,                 // the same as .WithReply
				Once:     true,                             // To not use it twice if true
			},
			{
				Pattern:  `DELETE "role_permission" WHERE`, // the same as .WithQuery()
				Response: Roleresp,                 // the same as .WithReply
				Once:     true,                             // To not use it twice if true
			},
			{
				Pattern:  `INSERT INTO "role_permission"`, // the same as .WithQuery()
				Response: Rolepermission,                 // the same as .WithReply
				Once:     true,                             // To not use it twice if true
			},

		})

		if assert.NoError(t, UpdateRole(c)) {
			spew.Dump(rec.Body)
			assert.Equal(t, http.StatusOK, rec.Code)
		}

		//test update role failed, role name already used
		req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(roleJSON))
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+ signToken)
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.SetPath("role/1")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  `SELECT * FROM "roles"  WHERE`, // the same as .WithQuery()
				Response: Rolerespfail2,                 // the same as .WithReply
				Once:     true,                             // To not use it twice if true
			},
			{
				Pattern:  `UPDATE "roles" SET`, // the same as .WithQuery()
				Response: Roleresp,                 // the same as .WithReply
				Once:     true,                             // To not use it twice if true
			},
			{
				Pattern:  `DELETE "role_permission" WHERE`, // the same as .WithQuery()
				Response: Roleresp,                 // the same as .WithReply
				Once:     true,                             // To not use it twice if true
			},
			{
				Pattern:  `INSERT INTO "role_permission"`, // the same as .WithQuery()
				Response: Rolepermission,                 // the same as .WithReply
				Once:     true,                             // To not use it twice if true
			},

		})

		if assert.NoError(t, UpdateRole(c)) {
			spew.Dump(rec.Body)
			assert.Equal(t, http.StatusOK, rec.Code)
		}


		//test update role failed, role not found
		req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(roleJSON))
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+ signToken)
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.SetPath("role/100")
		c.SetParamNames("id")
		c.SetParamValues("100")

		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  `SELECT * FROM "roles"  WHERE`, // the same as .WithQuery()
				Response: Rolerespfail,                 // the same as .WithReply
				Once:     true,                             // To not use it twice if true
			},

		})

		if assert.NoError(t, UpdateRole(c)) {

			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}




	})
}
