package config

import (
	"sync"

	"github.com/kelseyhightower/envconfig"
)

// Config is
type Config struct {
	AppEnv     string `envconfig:"APP_ENV" default:"development"`
	BcryptCost int    `envconfig:"BCRYPT_COST" default:"13"`

	CookieSecret  string `envconfig:"COOKIE_SECRET" default:"gEXYRjfN1gSXVuJXnI2x"`
	JwtKey        string `envconfig:"JWT_KEY" default:"85363D33AF8817405E9BD6650A8461AADE33E8D1CD7AF29C13A0C2CF025E35EA"`
	AdminUsername string `envconfig:"ADMIN_USER" default:"admin@data.com"`
	AdminPassword string `envconfig:"ADMIN_PASSWORD" default:"clearsoup"`

	DatabaseURL     string `envconfig:"DATABASE_URL" default:"postgres://postgres@localhost:5432/admin_dev?sslmode=disable"`
	TestDatabaseURL string `envconfig:"TEST_DATABASE_URL"`
}

const appPrefix = ""

var (
	instance *Config
	once     sync.Once
)

// GetInstance returns a Config pointer to retrieve environment variables
func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{}
		if err := envconfig.Process(appPrefix, instance); err != nil {
			panic(err)
		}
	})

	return instance
}
