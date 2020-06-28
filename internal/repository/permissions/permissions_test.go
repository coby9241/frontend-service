package permissions_test

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/coby9241/frontend-service/internal/models/permissions"
	. "github.com/coby9241/frontend-service/internal/repository/permissions"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
)

type PermRepoSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	repo Repository
}

func (r *PermRepoSuite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, r.mock, err = sqlmock.New()
	r.Assert().NoError(err)

	gormDB, err := gorm.Open("postgres", db)
	gormDB.LogMode(true)
	r.Assert().NoError(err)

	r.repo = NewPermissionsRepositoryImpl(gormDB)
}

func (r *PermRepoSuite) TestCreateNewRole() {
	// expect happy path
	r.T().Run("test success CreateNewRole", func(t *testing.T) {
		r.mock.ExpectBegin()
		r.mock.
			ExpectQuery(regexp.QuoteMeta(`INSERT INTO "roles" ("created_at","updated_at","name") VALUES ($1,$2,$3) RETURNING "roles"."id"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "tester").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(1))

		r.mock.
			ExpectQuery(`INSERT INTO "resources"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(1))

		r.mock.
			ExpectExec(`INSERT INTO "resource_role"`).
			WithArgs(1, 1, 1, 1).
			WillReturnResult(sqlmock.NewResult(1, 2))

		r.mock.ExpectCommit()

		role, err := r.repo.CreateNewRole([]*permissions.Resource{{ResourceName: "test"}}, "tester")
		r.Assert().NoError(err)
		// check role
		r.Equal("tester", role.Name)
		r.Equal(uint(1), role.ID)
		// check resource
		r.Equal(1, len(role.Resources))
		r.Equal(uint(1), role.Resources[0].ID)
		r.Equal("test", role.Resources[0].ResourceName)
	})

	// expect failure
	// r.T().Run("test failure CreateNewRole", func(t *testing.T) {
	// 	r.mock.
	// 		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL AND ((UID = $1)) ORDER BY "users"."id" ASC LIMIT 1`)).
	// 		WithArgs("test@data.com").
	// 		WillReturnError(errors.New("db error"))

	// 	usr, err := r.repo.GetUserByUID("test@data.com")
	// 	r.Assert().Error(err)
	// 	r.Equal(users.User{}, *usr)
	// })
}

func TestPermSuite(t *testing.T) {
	suite.Run(t, new(PermRepoSuite))
}
