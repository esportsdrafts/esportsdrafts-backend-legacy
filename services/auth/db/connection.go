
package db

import (
	"os"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// CreateDBHandler Creates a MySQL GORM driver
func CreateDBHandler(hostname string, user string, password string, dbname string) (*gorm.DB, error) {
	connParams := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, hostname, dbname)
	fmt.Fprintf(os.Stderr, "Conn params: %s\n", connParams)
	db, err := gorm.Open("mysql", connParams)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	db.DropTableIfExists(Account{})
	db.AutoMigrate(Account{})
	return db, nil
}
