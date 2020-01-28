package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/coby9241/frontend-service/internal/auth"
	"github.com/coby9241/frontend-service/internal/models/users"
	repo "github.com/coby9241/frontend-service/tests/mocks"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AuthSuite struct {
	suite.Suite
	auth       *AdminAuth
	testCtx    *admin.Context
	testCookie string
}

func (a *AuthSuite) SetupSuite() {
	// set up admin auth
	conf := &AdminAuthConfig{
		LoginPath:        "/testlogin",
		LogoutPath:       "/testlogout",
		CookieSecret:     "cookiesecret",
		SessionStoreName: "teststore",
		SessionStoreKey:  "uid",
	}

	a.auth = NewAdminAuth(conf, nil)

	// set test gin context
	context := &admin.Context{Admin: admin.New(&qor.Config{})}
	a.testCtx = context

	// set up test cookie
	a.setupTestCookie()
}

func (a *AuthSuite) setupTestCookie() {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Use(sessions.Sessions(a.auth.Sess.Name, a.auth.Sess.Store))
	r.GET("/set", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set(a.auth.Sess.Key, "user1")
		_ = session.Save()
		c.String(200, "ok")
	})

	res1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/set", nil)
	r.ServeHTTP(res1, req1)

	a.testCookie = res1.Header().Get("Set-Cookie")
}

func (a *AuthSuite) TestLoginURL() {
	a.Equal("/testlogin", a.auth.LoginURL(a.testCtx))
}

func (a *AuthSuite) TestLogoutURL() {
	a.Equal("/testlogout", a.auth.LogoutURL(a.testCtx))
}

func (a *AuthSuite) TestAdminAuthentication() {
	cases := []struct {
		name      string
		storeName string
		src       string
		gld       interface{}
		err       error
	}{
		{
			name:      "test auth get user success",
			storeName: "teststore",
			src:       a.testCookie,
			gld:       &users.User{},
		},
		{
			name:      "test auth get user cookie not found",
			storeName: "teststore",
			src:       "normalcookie",
		},
		{
			name:      "test auth get user not found from db",
			storeName: "teststore",
			src:       a.testCookie,
			err:       gorm.ErrRecordNotFound,
		},
		{
			name:      "test auth get user invalid cookie name",
			storeName: "teststore\x00",
			src:       a.testCookie,
		},
	}

	for _, tt := range cases {
		// prevent shadowing
		tc := tt
		a.T().Run(tc.name, func(t *testing.T) {
			// setup admin
			Admin := admin.New(&admin.AdminConfig{})

			// setup http req and resp
			res2 := httptest.NewRecorder()
			req2, _ := http.NewRequest("GET", "/get", nil)
			req2.Header.Set("Cookie", tc.src)
			ctx := Admin.NewContext(res2, req2)

			// setup mock
			mockRepo := &repo.UserRepository{}
			mockRepo.On("GetUserByUID", mock.Anything).Return(tc.gld, tc.err)
			a.auth.UserRepo = mockRepo

			// setup sess store name
			a.auth.Sess.Name = tc.storeName

			// test get current user and assert
			src := a.auth.GetCurrentUser(ctx)
			a.Equal(tc.gld, src)
		})
	}
}

func TestAuth(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}
