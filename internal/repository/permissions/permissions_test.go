package permissions_test

import (
	"database/sql"
	"errors"
	"fmt"
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

		// we make sure that all expectations were met
		if err := r.mock.ExpectationsWereMet(); err != nil {
			r.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	// expect failure
	r.T().Run("test failure CreateNewRole", func(t *testing.T) {
		r.mock.
			ExpectQuery(regexp.QuoteMeta(`INSERT INTO "roles" ("created_at","updated_at","name") VALUES ($1,$2,$3) RETURNING "roles"."id"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "tester").
			WillReturnError(errors.New("db error"))

		role, err := r.repo.CreateNewRole([]*permissions.Resource{{ResourceName: "test"}}, "tester")
		r.Assert().Error(err)
		r.Assert().Nil(role)

		// we make sure that all expectations were met
		if err := r.mock.ExpectationsWereMet(); err != nil {
			r.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func fixedFullRe(s string) string {
	return fmt.Sprintf("^%s$", regexp.QuoteMeta(s))
}

func (r *PermRepoSuite) TestSetPermissions() {
	// expect happy path
	r.T().Run("test success SetPermissions", func(t *testing.T) {
		r.mock.MatchExpectationsInOrder(false)
		r.mock.ExpectBegin()
		r.mock.
			ExpectExec(regexp.QuoteMeta(`UPDATE "resource_role" SET "can_create" = $1, "can_delete" = $2, "can_read" = $3, "can_update" = $4  WHERE (resource_id = $5 AND role_id = $6)`)).
			WithArgs(true, true, true, true, 1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		r.mock.ExpectCommit()

		err := r.repo.SetPermissions(1, 1, permissions.EnabledAttributes{
			CanCreate: true,
			CanRead:   true,
			CanUpdate: true,
			CanDelete: true,
		})
		r.Assert().NoError(err)

		// we make sure that all expectations were met
		if err := r.mock.ExpectationsWereMet(); err != nil {
			r.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	// expect failure path
	r.T().Run("test failure SetPermissions", func(t *testing.T) {
		mockErr := errors.New("test error")

		r.mock.MatchExpectationsInOrder(false)
		r.mock.ExpectBegin()
		r.mock.
			ExpectExec(regexp.QuoteMeta(`UPDATE "resource_role" SET "can_create" = $1, "can_delete" = $2, "can_read" = $3, "can_update" = $4  WHERE (resource_id = $5 AND role_id = $6)`)).
			WithArgs(true, true, true, true, 1, 1).
			WillReturnError(mockErr)

		r.mock.ExpectRollback()

		err := r.repo.SetPermissions(1, 1, permissions.EnabledAttributes{
			CanCreate: true,
			CanRead:   true,
			CanUpdate: true,
			CanDelete: true,
		})
		r.Assert().Equal(mockErr, err)

		// we make sure that all expectations were met
		if err := r.mock.ExpectationsWereMet(); err != nil {
			r.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestPermSuite(t *testing.T) {
	suite.Run(t, new(PermRepoSuite))
}
