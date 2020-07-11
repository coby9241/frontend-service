package migration

import (
	"github.com/coby9241/frontend-service/internal/db/migration/migrations"
	log "github.com/coby9241/frontend-service/internal/logger"
	"github.com/jinzhu/gorm"
	gormigrate "gopkg.in/gormigrate.v1"
)

// RunMigrations run all migrations in here
func RunMigrations(db *gorm.DB) error {
	// list of migrations to run
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{})

	// run init schema
	m.InitSchema(migrations.InitSchema)

	// run migration
	if err := m.Migrate(); err != nil {
		return err
	}

	log.GetInstance().Println("Migration ran successfully")

	return nil
}
