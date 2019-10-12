package db

import (
	"database/sql"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func createDatabase(hostname string, user string, password string, dbname string) {
	connParams := fmt.Sprintf("%s:%s@tcp(%s:3306)/", user, password, hostname)
	db, err := sql.Open("mysql", connParams)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbname)
	if err != nil {
		panic(err)
	}
}

// CreateDBHandler Creates a MySQL GORM driver
func CreateDBHandler(hostname string, user string, password string, dbname string) (*gorm.DB, error) {
	createDatabase(hostname, user, password, dbname)
	connParams := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, hostname, dbname)
	db, err := gorm.Open("mysql", connParams)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	db = db.AutoMigrate(Account{}, EmailVerificationCode{}, MFACode{}, MFAMethod{})
	return db, nil
}
