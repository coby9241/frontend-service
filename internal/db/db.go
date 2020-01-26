package db

import (
	"flag"
	"sync"

	"github.com/coby9241/frontend-service/internal/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // dialect imported to recognise it is using pg backend
	"github.com/qor/validations"
)

var instance *gorm.DB
var once sync.Once

// GetInstance is
func GetInstance() *gorm.DB {
	once.Do(func() {
		var err error
		var databaseURL string
		if flag.Lookup("test.v") == nil {
			databaseURL = config.GetInstance().DatabaseURL
		} else {
			databaseURL = config.GetInstance().TestDatabaseURL
		}

		if instance, err = gorm.Open("postgres", databaseURL); err != nil {
			panic(err)
		}

		instance.LogMode(true)
		validations.RegisterCallbacks(instance)
	})

	return instance
}
