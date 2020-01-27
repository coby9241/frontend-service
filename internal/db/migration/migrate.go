package migration

import (
	"time"

	"github.com/coby9241/frontend-service/internal/config"
	log "github.com/coby9241/frontend-service/internal/logger"
	"github.com/coby9241/frontend-service/internal/models/users"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	gormigrate "gopkg.in/gormigrate.v1"
)

// RunMigrations run all migrations in here
func RunMigrations(db *gorm.DB) error {
	// no migrations for now
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{})
	m.InitSchema(func(tx *gorm.DB) error {
		err := tx.AutoMigrate(
			&users.User{},
		).Error

		if err != nil {
			return err
		}

		// add admin user
		var pwd []byte
		if pwd, err = bcrypt.GenerateFromPassword([]byte(config.GetInstance().AdminPassword), bcrypt.DefaultCost); err != nil {
			return err
		}

		currTime := time.Now()

		usr := users.User{
			Provider:          "email",
			UID:               config.GetInstance().AdminUsername,
			PasswordHash:      string(pwd),
			UserID:            "admin",
			PasswordChangedAt: &currTime,
			CreatedAt:         currTime,
			UpdatedAt:         currTime,
			Role:              "admin",
		}

		if err = tx.Save(&usr).Error; err != nil {
			return err
		}

		// all other foreign keys...
		return nil
	})

	// run migration
	if err := m.Migrate(); err != nil {
		return err
	}

	log.GetInstance().Println("Migration ran successfully")

	return nil
}
