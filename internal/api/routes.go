package api

import (
	"net/http"

	v1 "github.com/coby9241/frontend-service/internal/api/v1"
	"github.com/coby9241/frontend-service/internal/auth"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes is
func RegisterAuthRoutes(r *gin.Engine, mux *http.ServeMux, adm *auth.AdminAuth) {
	g := r.Group("")
	g.Use(sessions.Sessions(adm.Sess.Name, adm.Sess.Store))
	{
		g.Any("/admin/*resources", gin.WrapH(mux))
		g.GET(adm.LoginPath, v1.GetLoginPage(adm))
		g.POST(adm.LoginPath, v1.LoginHandler(adm))
		g.GET(adm.LogoutPath, v1.GetLogout(adm))
	}
}
