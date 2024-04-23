package drivers

import (
	"fmt"
	"os"

	"github.com/base-al/base-core/helpers"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB connect to Postgress db
func ConnectPostgresDB() (*gorm.DB, error) {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Print("Error loading .env file")
	}
	//dbDriver := os.Getenv("DB_DRIVER")
	dbhost := os.Getenv("DB_HOST")
	dbusername := os.Getenv("DB_USERNAME")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	dbport := helpers.GetEnvAsInt("DB_PORT", 5432)

	dsn := fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%d sslmode=disable`, dbhost, dbusername, dbpassword, dbname, dbport)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
