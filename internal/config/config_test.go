package config

import (
	"sync"
	"testing"

	"github.com/coby9241/frontend-service/tests/utils"
	"github.com/stretchr/testify/assert"
)

func TestGoodConfig(t *testing.T) {
	cases := []struct {
		name    string
		src     []utils.EnvPair
		gld     *Config
		wantErr bool
	}{
		{
			name: "test load config",
			src: []utils.EnvPair{
				{Key: "APP_ENV", Value: "test"},
				{Key: "BCRYPT_COST", Value: "15"},
				{Key: "COOKIE_SECRET", Value: "cookiesecret"},
				{Key: "JWT_KEY", Value: "jwtkey"},
				{Key: "DATABASE_URL", Value: "postgres://postgres@localhost:5432/admin_dev?sslmode=disable"},
				{Key: "TEST_DATABASE_URL", Value: "postgres://postgres@localhost:5432/admin_dev?sslmode=disable"},
				{Key: "ADMIN_USER", Value: "user"},
				{Key: "ADMIN_PASSWORD", Value: "password"},
			},
			gld: &Config{
				AppEnv:          "test",
				BcryptCost:      15,
				CookieSecret:    "cookiesecret",
				JwtKey:          "jwtkey",
				DatabaseURL:     "postgres://postgres@localhost:5432/admin_dev?sslmode=disable",
				TestDatabaseURL: "postgres://postgres@localhost:5432/admin_dev?sslmode=disable",
				AdminUsername:   "user",
				AdminPassword:   "password",
			},
		},
		{
			name: "test load config with panic",
			src: []utils.EnvPair{
				{Key: "APP_ENV", Value: "test"},
				{Key: "BCRYPT_COST", Value: "XX"},
				{Key: "COOKIE_SECRET", Value: "cookiesecret"},
				{Key: "JWT_KEY", Value: "jwtkey"},
				{Key: "DATABASE_URL", Value: "postgres://cbloo@localhost:5432/admin_dev?sslmode=disable"},
				{Key: "TEST_DATABASE_URL", Value: "postgres://cbloo@localhost:5432/admin_dev?sslmode=disable"},
				{Key: "ADMIN_USER", Value: "user"},
				{Key: "ADMIN_PASSWORD", Value: "password"},
			},
			wantErr: true,
		},
	}

	for _, tt := range cases {
		// prevent shadowing
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			// reset sync once to retrigger load
			once = *new(sync.Once)
			// rescue panic if needed for loading
			defer func() {
				if err := recover(); (err != nil) != tc.wantErr {
					t.Errorf("failed test to load config. wantErr: %v, err: %v", tc.wantErr, err)
				}
			}()

			// set env vars
			resetEnv := utils.SetTestEnv(tc.src)
			defer resetEnv()
			// get conf
			conf := GetInstance()
			// assert conf is equal
			assert.Equal(t, tc.gld, conf)
		})
	}
}
