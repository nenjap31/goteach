package controllers

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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


func TestUser(t *testing.T){
	t.Run("test.v", func(t *testing.T) {
		e := echo.New()

		//test login success
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSON))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("login")

		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  `SELECT * FROM "users"  WHERE`, // the same as .WithQuery()
				Response: userResp2,                 // the same as .WithReply
				Once:     true,                             // To not use it twice if true
			},
			{
				Pattern:  `SELECT * FROM "users" WHERE`, // the same as .WithQuery()
				Response: userResp3,               // the same as .WithReply
				Once:     true,                               // To not use it twice if true
			},
			{
				Pattern:  `SELECT * FROM "users" WHERE`, // the same as .WithQuery()
				Response: UserResp,               // the same as .WithReply
				Once:     true,                               // To not use it twice if true
			},


		})

		if assert.NoError(t, LoginUser(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}

		//test login fail, user not found
		req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSONfail2))
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.SetPath("login")

		mocket.Catcher.Reset().NewMock().WithQuery("SELECT * FROM `users` WHERE").WithReply(UserResp)
		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  `SELECT * FROM "users"  WHERE`, // the same as .WithQuery()
				Response: userRespfail2,                 // the same as .WithReply
				Once:     true,                             // To not use it twice if true
			},
			{
				Pattern:  `SELECT * FROM "users" WHERE`, // the same as .WithQuery()
				Response: userResp3,               // the same as .WithReply
				Once:     true,                               // To not use it twice if true
			},
			{
				Pattern:  `SELECT * FROM "users" WHERE`, // the same as .WithQuery()
				Response: UserResp,               // the same as .WithReply
				Once:     true,                               // To not use it twice if true
			},


		})

		if assert.NoError(t, LoginUser(c)) {
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		}

		//test login fail, user not active
		req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSONfail))
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.SetPath("login")

		mocket.Catcher.Reset().NewMock().WithQuery("SELECT * FROM `users` WHERE").WithReply(UserResp)
		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  `SELECT * FROM "users"  WHERE`, // the same as .WithQuery()
				Response: userResp4,                 // the same as .WithReply
				Once:     true,                             // To not use it twice if true
			},
			{
				Pattern:  `SELECT * FROM "users" WHERE`, // the same as .WithQuery()
				Response: userResp3,               // the same as .WithReply
				Once:     true,                               // To not use it twice if true
			},
			{
				Pattern:  `SELECT * FROM "users" WHERE`, // the same as .WithQuery()
				Response: UserResp,               // the same as .WithReply
				Once:     true,                               // To not use it twice if true
			},


		})

		if assert.NoError(t, LoginUser(c)) {

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		}

		//test login fail, password wrong
		req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSONfail))
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.SetPath("login")

		mocket.Catcher.Reset().NewMock().WithQuery("SELECT * FROM `users` WHERE").WithReply(UserResp)
		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  `SELECT * FROM "users"  WHERE`, // the same as .WithQuery()
				Response: userResp2,                 // the same as .WithReply
				Once:     true,                             // To not use it twice if true
			},
			{
				Pattern:  `SELECT * FROM "users" WHERE`, // the same as .WithQuery()
				Response: userResp3,               // the same as .WithReply
				Once:     true,                               // To not use it twice if true
			},
			{
				Pattern:  `SELECT * FROM "users" WHERE`, // the same as .WithQuery()
				Response: UserResp,               // the same as .WithReply
				Once:     true,                               // To not use it twice if true
			},


		})

		if assert.NoError(t, LoginUser(c)) {
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		}


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

		//test get profile success
		req = httptest.NewRequest(http.MethodGet, "/user",nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+ signToken)
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.SetPath("/")
		h := middleware.JWTWithConfig(config.JwtConfig)(GetProfile)
		mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "users"`).WithReply(UserResp)
		// Assertions
		if assert.NoError(t, h(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}

		//test GetListUser success
		req = httptest.NewRequest(http.MethodGet, "/user",nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+ signToken)
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.SetPath("/list")
		h = middleware.JWTWithConfig(config.JwtConfig)(GetListUser)
		mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "users"`).WithReply(UserResp)
		// Assertions
		if assert.NoError(t, h(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}


		// test add user success
		req = httptest.NewRequest(http.MethodPost, "/user",strings.NewReader(userJSONAdd))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+ signToken)
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.SetPath("")

		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  `SELECT * FROM "users"  WHERE`,
				Response: userRespfail2,
				Once:     true,
			},
			{
				Pattern:  `SELECT * FROM "users"  WHERE`,
				Response: userRespfail2,
				Once:     true,
			},
			{
				Pattern:  `INSERT  INTO "users"`,
				Response: UserResp,
				Once:     true,
			},
		})
		// Assertions

		if assert.NoError(t, AddUser(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}

		// test add user failed, validation error
		req = httptest.NewRequest(http.MethodPost, "/user",strings.NewReader(userJSONAddFail))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+ signToken)
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.SetPath("")

		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  `SELECT * FROM "users"  WHERE`,
				Response: userRespfail2,
				Once:     true,
			},
			{
				Pattern:  `SELECT * FROM "users"  WHERE`,
				Response: userRespfail2,
				Once:     true,
			},
			{
				Pattern:  `INSERT  INTO "users"`,
				Response: UserResp,
				Once:     true,
			},
		})
		// Assertions

		if assert.NoError(t, AddUser(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}


		// test add user fail, username already used
		req = httptest.NewRequest(http.MethodPost, "/user",strings.NewReader(userJSONAdd))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+ signToken)
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.SetPath("")

		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  `SELECT * FROM "users"  WHERE`,
				Response: userRespfail,
				Once:     true,
			},
			{
				Pattern:  `SELECT * FROM "users"  WHERE`,
				Response: userRespfail2,
				Once:     true,
			},
			{
				Pattern:  `INSERT  INTO "users"`,
				Response: UserResp,
				Once:     true,
			},
		})
		// Assertions

		if assert.NoError(t, AddUser(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}


		// test add user fail, email already used
		req = httptest.NewRequest(http.MethodPost, "/user",strings.NewReader(userJSONAdd))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+ signToken)
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.SetPath("")

		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  `SELECT * FROM "users"  WHERE`,
				Response: userRespfail2,
				Once:     true,
			},
			{
				Pattern:  `SELECT * FROM "users"  WHERE`,
				Response: userRespfail,
				Once:     true,
			},
			{
				Pattern:  `INSERT  INTO "users"`,
				Response: UserResp,
				Once:     true,
			},
		})
		// Assertions

		if assert.NoError(t, AddUser(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}


		//test update user success
		req = httptest.NewRequest(http.MethodPut, "/user/1",strings.NewReader(userJSONupdate))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")
		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  `SELECT count(*) FROM "users"`,
				Response: UserCount,
				Once:     true,
			},
			{
				Pattern:  `SELECT * FROM "users"  WHERE`,
				Response: userRespfail2,
				Once:     true,
			},
			{
				Pattern:  `SELECT * FROM "users"  WHERE`,
				Response: userRespfail2,
				Once:     true,
			},
			{
				Pattern:  `SELECT * FROM "users"  WHERE`,
				Response: userRespfail2,
				Once:     true,
			},
			{
				Pattern:  `UPDATE "users" set`,
				Response: UserResp,
				Once:     true,
			},
		})
		// Assertions
		if assert.NoError(t, UpdateUser(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}

		//test update user success with new password
		req = httptest.NewRequest(http.MethodPut, "/user/1",strings.NewReader(userJSONupdate3))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")
		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  `SELECT count(*) FROM "users"`,
				Response: UserCount,
				Once:     true,
			},
			{
				Pattern:  `SELECT * FROM "users"  WHERE`,
				Response: userRespfail2,
				Once:     true,
			},
			{
				Pattern:  `SELECT * FROM "users"  WHERE`,
				Response: userRespfail2,
				Once:     true,
			},
			{
				Pattern:  `SELECT * FROM "users"  WHERE`,
				Response: userRespfail2,
				Once:     true,
			},
			{
				Pattern:  `UPDATE "users" set`,
				Response: UserResp,
				Once:     true,
			},
		})
		// Assertions
		if assert.NoError(t, UpdateUser(c)) {
			spew.Dump(rec.Body.String())
			assert.Equal(t, http.StatusOK, rec.Code)
		}

		//test update user failed, validation error username already exist
		req = httptest.NewRequest(http.MethodPut, "/user/2",strings.NewReader(userJSONupdate))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")
		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  `SELECT count(*) FROM "users"`,
				Response: UserCount,
				Once:     true,
			},
			{
				Pattern:  `SELECT * FROM "users"  WHERE`,
				Response: userRespfail2,
				Once:     true,
			},
			{
				Pattern:  `SELECT * FROM "users"  WHERE`,
				Response: userRespfail,
				Once:     true,
			},
			{
				Pattern:  `SELECT * FROM "users"  WHERE`,
				Response: userRespfail2,
				Once:     true,
			},
			{
				Pattern:  `UPDATE "users" set`,
				Response: UserResp,
				Once:     true,
			},
		})
		// Assertions
		if assert.NoError(t, UpdateUser(c)) {

			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}


		//test update user failed, validation error email already exist
		req = httptest.NewRequest(http.MethodPut, "/user/2",strings.NewReader(userJSONupdate))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")
		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  `SELECT count(*) FROM "users"`,
				Response: UserCount,
				Once:     true,
			},
			{
				Pattern:  `SELECT * FROM "users"  WHERE`,
				Response: userRespfail2,
				Once:     true,
			},
			{
				Pattern:  `SELECT * FROM "users"  WHERE`,
				Response: userRespfail2,
				Once:     true,
			},

			{
				Pattern:  `SELECT * FROM "users"  WHERE`,
				Response: userRespfail,
				Once:     true,
			},
			{
				Pattern:  `UPDATE "users" set`,
				Response: UserResp,
				Once:     true,
			},
		})
		// Assertions
		if assert.NoError(t, UpdateUser(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}

		//test update user failed, user not found
		req = httptest.NewRequest(http.MethodPut, "/user/2",strings.NewReader(userJSONupdate))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("100")
		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  `SELECT count(*) FROM "users"`,
				Response: userRespfail2,
				Once:     true,
			},
			{
				Pattern:  `SELECT * FROM "users"  WHERE`,
				Response: userRespfail2,
				Once:     true,
			},
			{
				Pattern:  `SELECT * FROM "users"  WHERE`,
				Response: userRespfail2,
				Once:     true,
			},
			{
				Pattern:  `SELECT * FROM "users"  WHERE`,
				Response: userRespfail2,
				Once:     true,
			},
			{
				Pattern:  `UPDATE "users" set`,
				Response: UserResp,
				Once:     true,
			},
		})
		// Assertions
		if assert.NoError(t, UpdateUser(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}



	})

}
