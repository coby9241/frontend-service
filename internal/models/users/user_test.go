package users_test

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/coby9241/frontend-service/internal/config"
	. "github.com/coby9241/frontend-service/internal/models/users"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
)

type UserSuite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
}

func (u *UserSuite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, u.mock, err = sqlmock.New()
	u.Assert().NoError(err)

	u.DB, err = gorm.Open("postgres", db)
	u.DB.LogMode(true)
	u.Assert().NoError(err)
}

func (u *UserSuite) TestCreateUser() {
	testUser := User{
		UserID:       "test",
		PasswordHash: "password",
	}

	// expect entire txn block for create
	// 1. begin
	// 2. insert
	// 3. commit
	u.mock.ExpectBegin()
	u.mock.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","uid","provider","password","role_id","user_id","password_changed_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "users"."id"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "", "", "password", 0, "test", nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(1))
	u.mock.ExpectCommit()

	if err := u.DB.Create(&testUser).Error; err != nil {
		u.T().Errorf("failed to create user: %v\n", err)
	}
}

func (u *UserSuite) TestGetResourceName() {
	u.Equal("user", User{}.GetResourceName())
}

func (u *UserSuite) TestDisplayName() {
	cases := []struct {
		name string
		src  *User
		gld  string
	}{
		{
			name: "test display with UID only",
			src: &User{
				UID: "x@y.com",
			},
			gld: "x@y.com",
		},
		{
			name: "test display with UserID and UID",
			src: &User{
				UID:    "x@y.com",
				UserID: "user",
			},
			gld: "user",
		},
	}

	for _, tt := range cases {
		tc := tt
		u.T().Run(tc.name, func(t *testing.T) {
			u.Equal(tc.gld, tc.src.DisplayName())
		})
	}
}

func (u *UserSuite) TestComparePassword() {
	cases := []struct {
		name     string
		src      *User
		password string
		isDiffPw bool
	}{
		{
			name: "test compare password same",
			src: &User{
				UID:          "x@y.com",
				PasswordHash: "$2a$13$lufjUdm8oJRBQWTHMJ96heLX5urtc5okRT8IdBGFpqz0DIKVfCnEu",
			},
			password: "password",
		},
		{
			name: "test compare password different",
			src: &User{
				UID:          "x@y.com",
				PasswordHash: "$2a$13$lufjUdm8oJRBQWTHMJ96heLX5urtc5okRT8IdBGFpqz0DIKVfCnEu",
			},
			password: "different",
			isDiffPw: true,
		},
	}

	for _, tt := range cases {
		tc := tt
		u.T().Run(tc.name, func(t *testing.T) {
			if err := tc.src.ComparePassword(tc.password); (err != nil) != tc.isDiffPw {
				u.T().Errorf("isDiffPw: %v, err: %v", tc.isDiffPw, err)
			}
		})
	}
}

func (u *UserSuite) TestGetJWTToken() {
	// set up test user
	timeNow := time.Now()
	testUser := &User{
		UID:               "x@y.com",
		PasswordChangedAt: &timeNow,
	}

	cases := []struct {
		name    string
		jwtKey  interface{}
		wantErr bool
	}{
		{
			name:   "test issue jwt token success",
			jwtKey: []byte(config.GetInstance().JwtKey),
		},
		{
			name:    "test issue jwt token with error",
			jwtKey:  123,
			wantErr: true,
		},
	}

	for _, tt := range cases {
		tc := tt
		u.T().Run(tc.name, func(t *testing.T) {
			if _, err := testUser.IssueJwtTokenSet(tc.jwtKey); (err != nil) != tc.wantErr {
				u.T().Errorf("failed to issue jwt token: %v\n", err)
			}
		})
	}
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}
