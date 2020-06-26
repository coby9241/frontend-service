package users_test

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

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
			WillReturnRows(sqlmock.NewRows([]string{"id", "uid"}).
				AddRow(1, "test@data.com"))

		usr, err := r.repo.GetUserByUID("test@data.com")
		r.Assert().NoError(err)
		r.Equal(users.User{
			Model: gorm.Model{
				ID: 1,
			},
			UID: "test@data.com",
		}, *usr)
	})

	// expect failure
	r.T().Run("test failure GetUserByUID", func(t *testing.T) {
		r.mock.
			ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL AND ((UID = $1)) ORDER BY "users"."id" ASC LIMIT 1`)).
			WithArgs("test@data.com").
			WillReturnError(errors.New("db error"))

		usr, err := r.repo.GetUserByUID("test@data.com")
		r.Assert().Error(err)
		r.Equal(users.User{}, *usr)
	})
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserRepoSuite))
}
