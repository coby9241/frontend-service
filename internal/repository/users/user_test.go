package users_test

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/coby9241/frontend-service/internal/models/permissions"
	"github.com/coby9241/frontend-service/internal/models/users"
	. "github.com/coby9241/frontend-service/internal/repository/users"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
)

type UserRepoSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	repo UserRepository
}

func (r *UserRepoSuite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, r.mock, err = sqlmock.New()
	r.Assert().NoError(err)

	gormDB, err := gorm.Open("postgres", db)
	gormDB.LogMode(true)
	r.Assert().NoError(err)

	r.repo = NewUserRepositoryImpl(gormDB)
}

func (r *UserRepoSuite) TestGetUserByUID() {
	// expect happy path
	r.T().Run("test success GetUserByUID", func(t *testing.T) {
		r.mock.
			ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL AND ((UID = $1)) ORDER BY "users"."id" ASC LIMIT 1`)).
			WithArgs("test@data.com").
			WillReturnRows(sqlmock.NewRows([]string{"id", "uid", "role_id"}).
				AddRow(1, "test@data.com", 1))

		r.mock.
			ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles"  WHERE ("id" = $1)`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
				AddRow(1, "tester"))

		usr, err := r.repo.GetUserByUID("test@data.com")
		r.Assert().NoError(err)
		r.Equal(users.User{
			Model: gorm.Model{
				ID: 1,
			},
			UID:    "test@data.com",
			RoleID: 1,
			Role: permissions.Role{
				ID:   1,
				Name: "tester",
			},
		}, *usr)
	})

	// expect failure in first query
	r.T().Run("test failure GetUserByUID", func(t *testing.T) {
		r.mock.
			ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL AND ((UID = $1)) ORDER BY "users"."id" ASC LIMIT 1`)).
			WithArgs("test@data.com").
			WillReturnError(errors.New("db error"))

		usr, err := r.repo.GetUserByUID("test@data.com")
		r.Assert().Error(err)
		r.Nil(usr)
	})

	r.T().Run("test failure GetUserByUID on relation loading", func(t *testing.T) {
		r.mock.
			ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL AND ((UID = $1)) ORDER BY "users"."id" ASC LIMIT 1`)).
			WithArgs("test@data.com").
			WillReturnRows(sqlmock.NewRows([]string{"id", "uid", "role_id"}).
				AddRow(1, "test@data.com", 1))
		r.mock.
			ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles"  WHERE ("id" = $1)`)).
			WithArgs(1).
			WillReturnError(errors.New("db error"))

		usr, err := r.repo.GetUserByUID("test@data.com")
		r.Assert().Error(err)
		r.Nil(usr)
	})
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserRepoSuite))
}
