package v1

import (
	"net/http"

	"github.com/coby9241/frontend-service/internal/auth"
	log "github.com/coby9241/frontend-service/internal/logger"
	"github.com/coby9241/frontend-service/internal/models/users"
	"github.com/coby9241/frontend-service/internal/response"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type (
	// LoginRequest is
	LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// SessionIssueResponse is
	SessionIssueResponse struct {
		Success  bool            `json:"success,omitempty"`
		TokenSet *users.TokenSet `json:"token_set,omitempty"`
	}
)

// GetLoginPage simply returns the login page
func GetLoginPage(auth *auth.AdminAuth) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		if sessions.Default(c).Get(auth.Sess.Key) != nil {
			c.Redirect(http.StatusSeeOther, "/admin")
			return
		}

		c.HTML(http.StatusOK, "login.html", gin.H{})
	})
}

// LoginHandler is
func LoginHandler(auth *auth.AdminAuth) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		realm := "Authorization Required"
		if auth.Realm == "" {
			realm = auth.Realm
		}

		// set auth realm
		c.Header("WWW-Authenticate", realm)

		session := sessions.Default(c)
		email := c.PostForm("email")
		password := c.PostForm("password")
		if email == "" || password == "" {
			response.RenderErrorPage(c, http.StatusUnauthorized, "missing username/password")
			return
		}

		i, err := auth.UserRepo.GetUserByUID(email)
		if err != nil {
			response.RenderErrorPage(c, http.StatusUnauthorized, "username/password incorrect")
			return
		}

		if err = i.ComparePassword(password); err != nil {
			response.RenderErrorPage(c, http.StatusUnauthorized, "username/password incorrect")
			return
		}

		session.Set(auth.Sess.Key, i.UID)
		if err := session.Save(); err != nil {
			log.GetInstance().WithError(err).Warn("Couldn't save session")
			c.Redirect(http.StatusInternalServerError, "unabled to login, please contact your administrator")
			return
		}

		c.Redirect(http.StatusSeeOther, "/admin")
	}

	return gin.HandlerFunc(fn)
}

// GetLogout allows the user to disconnect and redirects back to login page
func GetLogout(auth *auth.AdminAuth) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		session := sessions.Default(c)
		session.Delete(auth.Sess.Key)
		if err := session.Save(); err != nil {
			log.GetInstance().WithError(err).Warn("Couldn't save session")
		}

		c.Redirect(http.StatusSeeOther, auth.LoginPath)
	})
}
