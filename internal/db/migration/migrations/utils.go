package migrations

import "github.com/jinzhu/gorm"

// rollbackAndErr performs a rollback on the txn and returns an error
func rollbackAndErr(db *gorm.DB, err error) error {
	db.Rollback()
	return err
}
