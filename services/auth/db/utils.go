package db

import (
	"github.com/jinzhu/gorm"
)

// InTransaction denotes a function calling a GORM db
type InTransaction func(tx *gorm.DB) error

// DoInTransaction runs a InTransaction func in a transaction and rolls back if there
// is an error; otherwise commits the results
func DoInTransaction(fn InTransaction, db *gorm.DB) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	err := fn(tx)
	if err != nil {
		xerr := tx.Rollback().Error
		if xerr != nil {
			return xerr
		}
		return err
	}
	if err = tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
