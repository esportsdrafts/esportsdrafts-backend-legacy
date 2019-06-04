
package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// CreateDBHandler Creates a MySQL GORM driver
func CreateDBHandler(hostname string, user string, password string, dbname string) (*gorm.DB, error) {
	connParams := fmt.Sprintf("%s:%s@/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, dbname)
	db, err := gorm.Open(hostname, connParams)
	if err != nil {
		return nil, err
	}
	return db, nil
}
