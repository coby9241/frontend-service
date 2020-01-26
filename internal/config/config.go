package config

import (
	"sync"

	"github.com/kelseyhightower/envconfig"
)

// Config is
type Config struct {
	AppEnv     string `envconfig:"APP_ENV" default:"development"`
	BcryptCost int    `envconfig:"BCRYPT_COST" default:"13"`

	CookieSecret string `envconfig:"COOKIE_SECRET" default:"gEXYRjfN1gSXVuJXnI2x"`
	JwtKey       string `envconfig:"JWT_KEY" default:"gEXYRjfN1gSXVuJXnI2x"`

	DatabaseURL     string `envconfig:"DATABASE_URL" default:"postgres://cbloo@localhost:5432/admin_dev?sslmode=disable"`
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
