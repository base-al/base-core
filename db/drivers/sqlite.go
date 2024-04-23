package drivers

// ConnectSQLiteDB connect to SQLite db
import (
	"gorm.io/driver/sqlite" // Sqlite driver based on CGO
	"gorm.io/gorm"
)

func ConnectSQLiteDB() (*gorm.DB, error) {
	// github.com/mattn/go-sqlite3
	db, err := gorm.Open(sqlite.Open("db/base.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
