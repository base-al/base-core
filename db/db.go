package db

import (
	"os"

	drivers "github.com/base-al/base-core/db/drivers"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	// check if the DB_DRIVER is sqlite
	dbDriver := os.Getenv("DB_DRIVER")

	if dbDriver == "sqlite" {
		return drivers.ConnectSQLiteDB()
	} else if dbDriver == "mysql" {
		return drivers.ConnectMySQLDB()
	} else if dbDriver == "postgres" {
		return drivers.ConnectPostgresDB()
	} else {
		return nil, nil
	}
}
