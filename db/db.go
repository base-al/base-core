package db

import (
	"fmt"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := os.Getenv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

// ConnectDB connect to db
func ConnectDB() (*gorm.DB, error) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}

	dbhost := os.Getenv("DB_HOST")
	dbusername := os.Getenv("DB_USERNAME")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	dbport := getEnvAsInt("DB_PORT", 5432)

	dsn := fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%d sslmode=disable`, dbhost, dbusername, dbpassword, dbname, dbport)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
