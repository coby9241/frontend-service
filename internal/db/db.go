package db

import (
	"sync"

	"github.com/coby9241/frontend-service/internal/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // dialect imported to recognise it is using pg backend
	"github.com/qor/validations"
)

var instance *gorm.DB
var once sync.Once
var databaseURL string

// init will instantiate the uri to the one in the configuration and variable is sourced out to allow for
// dependency injection for testing
func init() {
	databaseURL = config.GetInstance().DatabaseURL
}

// GetInstance is
func GetInstance() *gorm.DB {
	once.Do(func() {
		var err error
		if instance, err = gorm.Open("postgres", databaseURL); err != nil {
			panic(err)
		}

		instance.LogMode(true)
		validations.RegisterCallbacks(instance)
	})

	return instance
}
