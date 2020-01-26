package auth

import (
	log "frontend-project/internal/logger"
	"frontend-project/internal/models/users"
	repo "frontend-project/internal/repository/users"

	"github.com/gin-contrib/sessions/cookie"
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor"
)

type (
	// AdminAuth is
	AdminAuth struct {
		UserRepo   repo.Repository
		LoginPath  string
		LogoutPath string
		Realm      string
		Sess       SessionStore
	}

	// AdminAuthConfig is
	AdminAuthConfig struct {
		LoginPath        string
		LogoutPath       string
		AuthRealm        string
		SessionStoreName string
		SessionStoreKey  string
		CookieSecret     string
	}

	// SessionStore is
	SessionStore struct {
		Name  string
		Key   string
		Store cookie.Store
	}
)

// NewAdminAuth is
func NewAdminAuth(conf *AdminAuthConfig, repo repo.Repository) *AdminAuth {
	return &AdminAuth{
		LoginPath:  conf.LoginPath,
		LogoutPath: conf.LogoutPath,
		Realm:      conf.AuthRealm,
		UserRepo:   repo,
		Sess: SessionStore{
			Name:  conf.SessionStoreName,
			Key:   conf.SessionStoreKey,
			Store: cookie.NewStore([]byte(conf.CookieSecret)),
		},
	}
}

// LoginURL is
func (a AdminAuth) LoginURL(c *admin.Context) string {
	return a.LoginPath
}

// LogoutURL is
func (a AdminAuth) LogoutURL(c *admin.Context) string {
	return a.LogoutPath
}

// GetCurrentUser is
func (a *AdminAuth) GetCurrentUser(c *admin.Context) qor.CurrentUser {
	session, err := a.Sess.Store.Get(c.Request, a.Sess.Name)
	if err != nil {
		return nil
	}

	var uid string
	if v, ok := session.Values[a.Sess.Key]; ok {
		uid = v.(string)
	} else {
		return nil
	}

	var id *users.User
	id, err = a.UserRepo.GetUserByUID(uid)
	if gorm.IsRecordNotFoundError(err) {
		log.GetInstance().WithError(err)
		return nil
	}

	return id
}
