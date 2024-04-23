package drivers

import (
	"fmt"
	"os"

	"github.com/base-al/base-core/helpers"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectMySQLDB() (*gorm.DB, error) {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Get database configuration from environment variables
	dbhost := os.Getenv("DB_HOST")
	dbusername := os.Getenv("DB_USERNAME")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	dbport := helpers.GetEnvAsInt("DB_PORT", 5432)

	// Connect to the database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", dbusername, dbpassword, dbhost, dbport)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Create the database if it doesn't exist
	result := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbname))
	if result.Error != nil {
		return nil, result.Error
	}

	// Connect to the specified database
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", dbusername, dbpassword, dbhost, dbport, dbname)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
