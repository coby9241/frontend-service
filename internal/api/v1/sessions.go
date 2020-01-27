package v1

import (
	"errors"
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

	// LoginForm is
	LoginForm struct {
		Email    string `form:"email"`
		Password string `form:"password"`
	}
)

var (
	// ErrorUserPasswordIncorrect is
	ErrorUserPasswordIncorrect = errors.New("username/password incorrect")
	// ErrorInternalServer is
	ErrorInternalServer = errors.New("unable to login, please contact your administrator")
)

// GetLoginPage simply returns the login page
func GetLoginPage(auth *auth.AdminAuth) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// redirect to admin page directly if cookie is found
		if sessions.Default(c).Get(auth.Sess.Key) != nil {
			c.Redirect(http.StatusSeeOther, "/admin")
			return
		}

		// render login page
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

		// get default session
		session := sessions.Default(c)

		// bind inputs from from to struct
		var l LoginForm
		if err := c.ShouldBind(&l); err != nil {
			log.GetInstance().WithError(err).Warn("Couldn't bind form contents to login request struct")
			response.RenderErrorPage(c, http.StatusInternalServerError, ErrorInternalServer)
			return
		}

		// check user by UID aka email
		i, err := auth.UserRepo.GetUserByUID(l.Email)
		if err != nil {
			response.RenderErrorPage(c, http.StatusUnauthorized, ErrorUserPasswordIncorrect)
			return
		}

		// validate password
		if err = i.ComparePassword(l.Password); err != nil {
			response.RenderErrorPage(c, http.StatusUnauthorized, ErrorUserPasswordIncorrect)
			return
		}

		// set session cookie
		session.Set(auth.Sess.Key, i.UID)
		if err = session.Save(); err != nil {
			log.GetInstance().WithError(err).Warn("Couldn't save session")
			response.RenderErrorPage(c, http.StatusInternalServerError, ErrorInternalServer)
			return
		}

		// redirect to admin page
		c.Redirect(http.StatusSeeOther, "/admin")
	}

	return gin.HandlerFunc(fn)
}

// GetLogout allows the user to disconnect and redirects back to login page
func GetLogout(auth *auth.AdminAuth) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// delete session key from cookie
		session := sessions.Default(c)
		session.Delete(auth.Sess.Key)
		if err := session.Save(); err != nil {
			log.GetInstance().WithError(err).Warn("Couldn't save session")
		}

		// redirect back to login page
		c.Redirect(http.StatusSeeOther, auth.LoginPath)
	})
}
