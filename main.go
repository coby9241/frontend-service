package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coby9241/frontend-service/internal/api"
	"github.com/coby9241/frontend-service/internal/auth"
	"github.com/coby9241/frontend-service/internal/bindatafs"
	"github.com/coby9241/frontend-service/internal/config"
	"github.com/coby9241/frontend-service/internal/db"
	"github.com/coby9241/frontend-service/internal/db/migration"
	"github.com/coby9241/frontend-service/internal/encryptor"
	log "github.com/coby9241/frontend-service/internal/logger"
	"github.com/coby9241/frontend-service/internal/models/permissions"
	"github.com/coby9241/frontend-service/internal/models/users"
	"github.com/coby9241/frontend-service/internal/rbac"
	permRepo "github.com/coby9241/frontend-service/internal/repository/permissions"
	userRepo "github.com/coby9241/frontend-service/internal/repository/users"
	"github.com/gin-gonic/gin"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	"github.com/qor/validations"
)

func main() {
	// set up database and run migrations
	DB := db.GetInstance()
	if err := migration.RunMigrations(DB); err != nil {
		panic(fmt.Errorf("failed to run migrations due to the following error: %v", err))
	}

	// set up qor admin interface
	admAuthConf := &auth.AdminAuthConfig{
		LoginPath:        "/login",
		LogoutPath:       "/logout",
		SessionStoreName: "admsession",
		SessionStoreKey:  "uid",
		CookieSecret:     config.GetInstance().CookieSecret,
	}
	admAuth := auth.NewAdminAuth(admAuthConf, userRepo.NewUserRepositoryImpl(DB))
	adm := admin.New(&admin.AdminConfig{
		DB:   DB,
		Auth: admAuth,
	})

	// get rbac repo
	permissionsRepo := permRepo.NewPermissionsRepositoryImpl(DB)
	_ = rbac.Load(permissionsRepo)
	// set resources in qor admin
	addUserResources(adm, permissionsRepo)

	router := gin.New()
	mountAssetFiles(router)
	initializeRoutes(router, adm, admAuth)

	// run router and wait for termination
	listenAndServe(router)
}

func listenAndServe(r *gin.Engine) {
	// init server on port 8082
	srv := &http.Server{
		Addr:    ":8082",
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.GetInstance().Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.GetInstance().Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.GetInstance().Fatal("Server Shutdown: ", err)
	}

	log.GetInstance().Println("Server exiting")
}

func initializeRoutes(r *gin.Engine, adm *admin.Admin, admAuth *auth.AdminAuth) {
	mux := mountAdmin(adm, admAuth)
	api.RegisterAuthRoutes(r, mux, admAuth)
}

func mountAdmin(adm *admin.Admin, admAuth *auth.AdminAuth) *http.ServeMux {
	mux := http.NewServeMux()
	adm.MountTo("/admin", mux)
	return mux
}

func mountAssetFiles(r *gin.Engine) {
	// mount template asset files
	lfs := bindatafs.AssetFS.NameSpace("login")
	err := lfs.RegisterPath("templates/")
	if err != nil {
		log.GetInstance().WithError(err).Fatal("Unable to register template folder for static pages in admin")
	}

	// set html template files
	logintpl, err := lfs.Asset("login.html")
	if err != nil {
		log.GetInstance().WithError(err).Fatal("Unable to find HTML template for login page in admin")
	}

	errtpl, err := lfs.Asset("error.tpl")
	if err != nil {
		log.GetInstance().WithError(err).Fatal("Unable to find HTML template for error page in admin")
	}

	// set html templates
	tpl := template.Must(template.New("login.html").Parse(string(logintpl)))
	tpl = template.Must(tpl.New("error.tpl").Parse(string(errtpl)))

	r.SetHTMLTemplate(tpl)

	// load css file
	r.StaticFile("main.css", "./templates/main.css")
}

func addUserResources(adm *admin.Admin, repo permRepo.Repository) {
	// get permissions for user resource
	userPermissions, err := rbac.ResourceRBAC(users.User{}.GetResourceName(), repo)
	if err != nil {
		panic(err)
	}

	user := adm.AddResource(&users.User{}, &admin.Config{
		Menu:       []string{"User Management"},
		Permission: userPermissions,
	})
	user.IndexAttrs("-PasswordHash", "-RoleID")
	user.Meta(&admin.Meta{
		Name: "PasswordHash",
		Type: "password",
		Setter: func(resource interface{}, metaValue *resource.MetaValue, context *qor.Context) {
			values := metaValue.Value.([]string)
			if len(values) > 0 {
				if np := values[0]; np != "" {
					pwd, err := encryptor.GetInstance().Digest(np)
					if err != nil {
						context.DB.AddError(validations.NewError(user, "Password", "Can't encrypt password")) // nolint: gosec,errcheck
						return
					}
					u := resource.(*users.User)
					u.PasswordHash = pwd
				}
			}
		},
	})
	user.Meta(&admin.Meta{
		Name: "PasswordChangedAt",
		Type: "datetime",
		Setter: func(resource interface{}, metaValue *resource.MetaValue, context *qor.Context) {
			u := resource.(*users.User)
			now := time.Now()
			u.PasswordChangedAt = &now
		},
	})
	user.Meta(&admin.Meta{
		Name: "Role",
		Type: "select_one",
		Config: &admin.SelectOneConfig{
			Collection: func(_ interface{}, context *admin.Context) (options [][]string) {
				var roles []permissions.Role
				context.DB.Find(&roles)

				for _, n := range roles {
					idStr := fmt.Sprintf("%d", n.ID)
					var option = []string{idStr, n.Name}
					options = append(options, option)
				}

				return options
			},
		},
		Valuer: func(user interface{}, ctx *qor.Context) interface{} {
			// for new user
			if user.(*users.User).ID == 0 {
				return true
			}

			// load relation and get username
			var role permissions.Role
			if err = ctx.DB.Model(user.(*users.User)).Related(&role).Error; err != nil {
				ctx.AddError(errors.New("failed to find role for user"))
				return nil
			}
			return role.Name
		},
		Setter: func(resource interface{}, metaValue *resource.MetaValue, context *qor.Context) {
			var role permissions.Role
			roleID := metaValue.Value.([]string)[0]
			context.DB.Where("id = ?", roleID).First(&role)

			user := resource.(*users.User)
			user.RoleID = role.ID
			metaValue.Value.([]string)[0] = role.Name
		},
	})

	user.ShowAttrs("Provider", "UID", "UserID", "Role")
	user.NewAttrs("Provider", "UID", "PasswordHash", "UserID", "Role")
	user.EditAttrs("Provider", "UID", "PasswordHash", "UserID", "Role")
}
