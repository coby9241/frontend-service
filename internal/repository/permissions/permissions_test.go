package permissions_test

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

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
		r.mock.MatchExpectationsInOrder(false)
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
		r.mock.MatchExpectationsInOrder(false)
		r.mock.ExpectBegin()
		r.mock.
			ExpectQuery(regexp.QuoteMeta(`INSERT INTO "roles" ("created_at","updated_at","name") VALUES ($1,$2,$3) RETURNING "roles"."id"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "tester").
			WillReturnError(errors.New("db error"))

		r.mock.ExpectRollback()

		role, err := r.repo.CreateNewRole([]*permissions.Resource{{ResourceName: "test"}}, "tester")
		r.Assert().Error(err)
		r.Assert().Nil(role)

		// we make sure that all expectations were met
		if err := r.mock.ExpectationsWereMet(); err != nil {
			r.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})
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

func (r *PermRepoSuite) TestGetRoles() {
	// fix time for testing
	testTime := time.Now()
	testError := errors.New("GetRoles error")

	cases := []struct {
		name       string
		mockExpect func(mock sqlmock.Sqlmock)
		golden     []permissions.Role
		wantErr    bool
	}{
		{
			name: "test success GetRoles",
			mockExpect: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles"`)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name"}).
						AddRow(1, testTime, testTime, "admin").
						AddRow(2, testTime, testTime, "editor"))
			},
			golden: []permissions.Role{
				{
					Name:      "admin",
					ID:        1,
					CreatedAt: testTime,
					UpdatedAt: testTime,
				},
				{
					Name:      "editor",
					ID:        2,
					CreatedAt: testTime,
					UpdatedAt: testTime,
				},
			},
		},
		{
			name: "test success GetRoles empty",
			mockExpect: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles"`)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name"}))
			},
			golden: []permissions.Role{},
		},
		{
			name: "test failure GetRoles",
			mockExpect: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles"`)).
					WillReturnError(testError)
			},
			golden:  nil,
			wantErr: true,
		},
	}

	for _, tt := range cases {
		tc := tt
		r.T().Run(tc.name, func(t *testing.T) {
			tc.mockExpect(r.mock)

			// we make sure that all expectations were met
			defer func() {
				if err := r.mock.ExpectationsWereMet(); err != nil {
					r.T().Errorf("there were unfulfilled expectations: %s", err)
				}
			}()

			roles, err := r.repo.GetRoles()
			if (err != nil) != tc.wantErr {
				t.Errorf("failed test to get roles. wantErr: %v, err: %v", tc.wantErr, err)
			}

			r.Equal(tc.golden, roles)
		})
	}
}

func TestPermSuite(t *testing.T) {
	suite.Run(t, new(PermRepoSuite))
}
